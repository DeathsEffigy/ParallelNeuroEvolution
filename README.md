## What is PNE?
Parallel Neuro-Evolution (PNE) is an experimental artificial intelligence that works very much unlike any other models currently out there. Basically, PNE does not use matrices at all. Instead, it emulates (with some degree of accuracy) a parallel activation paradigm where neurons either fire or don't fire, but they do not do so in sync. Consequently, the model will activate (and re-activate) a wide range of mechanical neurons (output neurons), that are ranked by frequency of activation and reactively inhibited. The neural model then evolves using forward propagation.

## What's next for PNE?
Right now, PNE works quite well. There is, however, a caveat: Since the model uses Goroutines to emulate a parallel activation paradigm, the model can get quite intense (and slow). Learning rate is perfectly fine and rather fast, but processing takes a while. This could be combatted by writing a pne_Brain.go container that allows construction of several pne_Circuit.go structs that perform quicker, much more focused tasks. This will be implemented some time in the future.

## Should I use PNE?
If you're asking this, the answer is probably no. This isn't for software development (let alone production releases). This is mostly for people interested in computational neurosciences.