/*
Package yfast implements a y-fast trie.  Instead of a red-black BBST
for the leaves, this implementation uses a simple ordered list.  This
package should have searches that are as performant as the x-fast
trie while having faster inserts/deletes and linear space consumption.

Performance characteristics:
Space: O(n)
Get: O(log log M)
Search: O(log log M)
Insert: O(log log M)
Delete: O(log log M)

where n is the number of items in the trie and M is the size of the
universe, ie, 2^m where m is the number of bits in the specified key
size.
*/

package yfast
