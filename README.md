# ancli-cli
Exponential-ish backoff for your notes!

`ancli` (Anki + CLI) is a CLI tool to make spaced repetition of your notes easy as easy as `ancli test`.

## Install

## What is ancli?

`ancli` is a CLI tool that attempts to reduce friction for users wanting to used space repetition to better remember their notes. `ancli` scrapes a user-defined directory, likely a notes directory, looking for regex matches on a decks configured symbol. Notes are then quizzed on an interval decided by the configured algorithm for the deck.

Features Include:
  - multiple decks
  - support multiple spaced repetition algorithms, including [SM2](https://www.supermemo.com/en/archives1990-2015/english/ol/sm2)
  - saving deck progress to file system

## Creating a Card

`ancli` is intended to be low overhead. With a few small changes to your standard note taking habits, you can begin creating question and answer cards for spaced repetition learning.

### Standard Style

To create an `ancli` card using the standard style, just place `---` anywhere in the file.

```
What kind of band plays snappy music?
---
A rubber band

#ancli-jokes
```

The `#ancli-jokes` tag is used to denote this card belongs to the jokes deck. 

### Inline Style

Similarly, the inline style can be used for answers that fall inline with the question. The answer part of the card will be overwritten with `_`'s during a quiz.

```
The ~trunk~ is where an elephant stores its suitcase.
#ancli-jokes
```

