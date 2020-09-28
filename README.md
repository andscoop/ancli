# ancli-cli
Exponential-ish backoff for your notes!

`ancli` (Anki + CLI) is a CLI tool to make spaced repetition of your notes easy as easy as `ancli test`.

## Install

## What is ancli?

`ancli` is a CLI tool that attempts to reduce friction for users wanting to used space repetition for increased knowledge retention. 

Just point `ancli` to your notes directory, give it a regex identifier for notes you want to be quizzed on, and begin testing your knowledge.

Features Include:
  - multiple decks
  - support multiple spaced repetition algorithms, including [SM2](https://www.supermemo.com/en/archives1990-2015/english/ol/sm2)
  - saving deck progress to file system

## Get Started With Decks

Decks are groupings of notes, or `cards` in `ancli`, that you want to be quizzed on. To get started and try out `ancli`, use the following command to create a deck of jokes you want to remember using the `example/` directory in this repo.

```
ancli decks create jokes "#ancli-jokes" -f examples/
```

This command walks the directory looking for the `#ancli-jokes` tag and building an index of any files with that tag.

You can then see a list of all the decks you have built.

```
ancli decks
```

`#ancli-jokes` is an arbitrary tag we have added to the notes text files to denote that we want to be quizzed on this note.

## Creating a Card

A card is a single unit of knowledge you want to remember. In practical terms, it is a note you have deemed important. 

An `ancli` card can be just any old txt or markdown file, but it's more fun if you use the below syntax to differentiate between the question and answer in the card.

### Standard Syntax

To create an `ancli` card using the standard syntax, just place `---` anywhere in the file.

```
What kind of band plays snappy music?
---
A rubber band

#ancli-jokes
```

### Inline Syntax

Similarly, the inline style can be used for answers that fall inline with the question. The answer part of the card will be overwritten with `_`'s during a quiz.

```
The ~trunk~ is where an elephant stores its suitcase.
#ancli-jokes
```
