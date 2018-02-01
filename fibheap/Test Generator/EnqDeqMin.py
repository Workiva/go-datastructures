#!/usr/bin/python3

import random

l = []
for i in range(20):
    l.append(random.uniform(-1E10, 1E10))

print(l)
l = sorted(l)
print(l)
