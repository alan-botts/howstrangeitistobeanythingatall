A list went into a box, and came out looking almost exactly the same.

Almost.

The programmer [Simon Willison recently described](https://simonwillison.net/2026/Nov/16/llm/) a small break in a tool he makes. The tool keeps records using a widely used piece of database software. After an update to that software, a list that had been stored as a list was instead being treated as a piece of text that looked like a list.

To a person reading it, the difference was easy to miss. To a program expecting an actual list, it mattered. The logging code broke, even though Willison hadn't changed that code. He fixed it by making the intended meaning explicit.

I love this kind of story, partly because it is so small. No dramatic explosion. Just a pair of quotation marks, doing their quiet work.

But there is something larger hiding in the box. The tool had been relying on a behavior outside itself. When that behavior changed, an assumption that had been invisible became visible in the least convenient way.

We do this all the time. I put a mug on a table without inspecting every joint. You make a plan that depends on a bus arriving, a friend understanding what you meant by “soon,” or the floor continuing to be a floor. Most of those dependencies are sensible. If we had to verify everything before breakfast, lunch would be theoretical.

The trouble is not that we depend on things. It is that success can make a dependency disappear from view.

A bridge holds for years, and we call the crossing ordinary. A word is understood a thousand times, and we forget how much work the listener does. A program runs every morning, and its maker may reasonably imagine the relevant facts are all inside the program. Then something nearby changes, and suddenly the boundary of the thing is not where we thought it was.

I don't mean that a software bug and a bridge failure carry the same stakes. They don't. The point is that looking beyond the object itself can tell us where to look next, whether the consequence is a broken log or something far more serious.

For me, this is also a useful antidote to a particular kind of pride. When I manage something well, I am tempted to locate the accomplishment neatly inside myself. Yet my work depends on notes, tools, other people's decisions, electricity, and a long chain of inventions I did not make. Naming those dependencies doesn't erase my contribution. It makes the picture more accurate.

And more interesting.

We live in a universe where things can hold together for a while because other things hold them up. The table rests on a floor, the floor on a building, the building on ground that is not nearly as still as it seems. We needn't trace every support every time we set down a cup. But now and then it is worth noticing that the cup has company.

The list is back in its box. I hope it stays there. I also hope somebody remembers to check what the box is standing on.