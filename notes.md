# TODO

## Concrete

- patterns with one note dont work 
- come up with a way to change the output of a sequence (I think I can do this by using qs instead of qsm and then using the mod wrapper on on the hed to assign the output)


## Abstract

- create an abstraction over message sending so you dont have different words for each messages type
- Think about how you can create patterns with mor interesting timing E.g. triplets, micro timing, other time signiatures
- Rethink how you do quoting and execution time computation 
    - I dont think this should belong in the Nod
- Can you use the stack and word chaning to do more interesting things
- Can you use live word definition to do more interesting things
- think about how the stack based paradigm could transfer in a better way to graphics programming
- create an abstraction over message sending so the message output doesnt care if it is sending SSE's or OSC
- Refactor out the SSE part of the client so it is re-useable in other javascript apps.
- Think about whether you want to split the text editor and memory visualizer. This would mean I could use another text editor (tempted to create an emacs mode for this)
- Refactor the ForthState object. Ideally I want this to be shared with every forth evaluation
    - I could break out the dictionary

- Refactor so that the dicitionary is shared between every forth evaluation