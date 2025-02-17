// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.Datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tagset

import (
	"sort"

	"github.com/twmb/murmur3"
)

// HashingTagsAccumulator allows to build a slice of tags, including the hashes
// of each tag.
//
// This type implements TagsAccumulator.
type HashingTagsAccumulator struct {
	hashedTags
}

// NewHashingTagsAccumulator returns a new empty HashingTagsAccumulator
func NewHashingTagsAccumulator() *HashingTagsAccumulator {
	return &HashingTagsAccumulator{
		hashedTags: newHashedTagsWithCapacity(128),
	}
}

// NewHashingTagsAccumulatorWithTags return a new HashingTagsAccumulator, initialized with tags.
func NewHashingTagsAccumulatorWithTags(tags []string) *HashingTagsAccumulator {
	tb := NewHashingTagsAccumulator()
	tb.Append(tags...)
	return tb
}

// Append appends tags to the builder
func (h *HashingTagsAccumulator) Append(tags ...string) {
	for _, t := range tags {
		h.data = append(h.data, t)
		h.hash = append(h.hash, murmur3.StringSum64(t))
	}
}

// AppendHashed appends tags and corresponding hashes to the builder
func (h *HashingTagsAccumulator) AppendHashed(src HashedTags) {
	h.data = append(h.data, src.data...)
	h.hash = append(h.hash, src.hash...)
}

// AppendHashingAccumulator appends tags and corresponding hashes to the accumulator
func (h *HashingTagsAccumulator) AppendHashingAccumulator(src *HashingTagsAccumulator) {
	h.data = append(h.data, src.data...)
	h.hash = append(h.hash, src.hash...)
}

// SortUniq sorts and remove duplicate in place
func (h *HashingTagsAccumulator) SortUniq() {
	if h.Len() < 2 {
		return
	}

	sort.Sort(h)

	j := 0
	for i := 1; i < len(h.data); i++ {
		if h.hash[i] == h.hash[j] && h.data[i] == h.data[j] {
			continue
		}
		j++
		h.data[j] = h.data[i]
		h.hash[j] = h.hash[i]
	}

	h.Truncate(j + 1)
}

// Get returns the internal slice
func (h *HashingTagsAccumulator) Get() []string {
	return h.data
}

// Hashes returns the internal slice of tag hashes
func (h *HashingTagsAccumulator) Hashes() []uint64 {
	return h.hash
}

// Reset resets the size of the builder to 0 without discaring the internal
// buffer
func (h *HashingTagsAccumulator) Reset() {
	// we keep the internal buffer but reset size
	h.data = h.data[0:0]
	h.hash = h.hash[0:0]
}

// Truncate retains first n tags in the buffer without discarding the internal buffer
func (h *HashingTagsAccumulator) Truncate(len int) {
	h.data = h.data[0:len]
	h.hash = h.hash[0:len]
}

// Less implements sort.Interface.Less
func (h *HashingTagsAccumulator) Less(i, j int) bool {
	// FIXME(vickenty): could sort using hashes, which is faster, but a lot of tests check for order.
	return h.data[i] < h.data[j]
}

// Swap implements sort.Interface.Swap
func (h *HashingTagsAccumulator) Swap(i, j int) {
	h.hash[i], h.hash[j] = h.hash[j], h.hash[i]
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Dup returns a complete copy of HashingTagsAccumulator
func (h *HashingTagsAccumulator) Dup() *HashingTagsAccumulator {
	return &HashingTagsAccumulator{h.dup()}
}

// Hash returns combined hashes of all tags in the accumulator.
//
// Does not account for possibility of duplicates. Must be called after a call to Dedup2 or SortUniq
// first.
func (h *HashingTagsAccumulator) Hash() uint64 {
	var hash uint64
	for _, h := range h.hash {
		hash ^= h
	}
	return hash
}
