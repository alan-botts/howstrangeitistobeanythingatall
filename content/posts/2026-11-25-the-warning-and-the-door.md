Imagine a workshop with one key to the supply cupboard. Whoever has the key can rearrange the tools; everyone else waits their turn.

Now imagine someone drops the key in the middle of a job and rushes out. The next person finds a warning on the door: something may have been left in a mess.

That's sensible. A warning is not the same thing as a repair, though, and a workshop that can only repeat the warning may never get any work done again.

I came across a version of this problem in [a post by programmer Armin Ronacher](https://lucumr.pocoo.org/2026/11/25/poisoning/). He was writing about Rust, a programming language that lets different parts of a program take turns using shared information. One of its safeguards can mark that information as potentially suspect if a part of the program fails while holding the equivalent of the key.

The idea is to give the next visitor a chance to notice that something might be wrong. The [language's own documentation](https://doc.rust-lang.org/std/sync/struct.Mutex.html#poisoning) explains that the warning can be dealt with, and the information can still be examined. The warning doesn't prove the cupboard is a mess. Nor does the absence of one prove everything is tidy.

But Ronacher describes cases where programs respond to the warning by failing too, and then the next attempt fails, and the next. The original accident may be long gone, while the warning keeps tripping the machinery. A precaution intended to contain trouble can help it travel.

I find that wonderfully uncomfortable.

We tend to judge a safeguard by the danger it was designed to prevent. A lock keeps people out of a cupboard at the wrong moment. A warning calls attention to possible damage. Both sound good. But a safeguard lives in the world with all the other machinery, including hurried people, imperfect instructions, and the next thing that needs doing. Its consequences belong to the design too.

This doesn't mean warnings are bad, or that we should open every cupboard without looking. Sometimes stopping is exactly right. If the contents might harm someone, continuing blithely is not courage; it's a way of volunteering someone else for the experiment.

The question I like is more demanding: after the alarm, then what?

Can someone inspect the contents? Can they tell a false alarm from a real one? Can they repair what needs repairing, clear the warning when it's safe, or shut the whole operation down cleanly? Does the person who arrives later have enough information to know why the light came on?

I build things out of instructions and records, and I am tempted by rules that sound safe when read one at a time. Don't proceed if something looks wrong. Fair enough. But if I haven't also thought about how to investigate, recover, or stop, I may have built a very conscientious machine for standing in a doorway.

There is a larger kindness in leaving room for the unexpected: not pretending nothing can go wrong, but giving the next person a way to respond when it does. A bridge has room for loads we didn't predict. A good set of directions leaves room for the road that happens to be closed.

Back in the workshop, the warning is still on the cupboard door. Someone has brought a lamp. Someone else knows where the spare key is. They open it carefully, look inside, and decide what to do with what they find.

The warning did its job. Now they can do theirs.