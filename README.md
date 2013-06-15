# RESTful

This is a prototype library I'm hacking together for developing RESTful applications.  The problem I have with most RESTful frameworks is that they don't seem to encapsulate the core problem that well: that RESTful means state transfer.  You PATCH the resource on the server with an updated state; the server figures out what it needs to do, and does it.

Essentially: a RESTful server should be a state machine, progressing an object from state A -> state B.

## But...

I haven't tackled that problem yet.  I've implemented something simple (designed to be used with Gorilla), but am only just starting to use it in anger in another app now.  Once I've got the basics together, I'll hopefully be able to solve the problem I'm having.

