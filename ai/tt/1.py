# coding: utf-8
get_ipython().run_line_magic('load', '1.py')
# %load 1.py
import torch
a = open('../../../makemore/names.txt','r').read().splitlines()
print(a)
b = {}
for w in a:
    chs = ['<S>'] + list(w) + ['<E>']
    for ch1, ch2 in zip(chs, chs[1:]):
        bigram = (ch1, ch2)
        b[bigram] = b.get(bigram, 0) + 1
sorted(b.items(), key = lambda kv: kv[1])
import torch
get_ipython().run_line_magic('save', '1')
p = torch.zeros((3,5))
p
