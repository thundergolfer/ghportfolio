# ghportfolio
A minimal CLI in Golang for catching up on the activity and health of your public Github projects

## Usage

#### `totals`

A quick way to tally up *Stars* and *Forks* across all your public repositories.

```
$ ghportfolio totals

Portfolio Stars: 312  Portfolio Forks: 76
```

#### `list`

Produces a table of all your public repositories, highlighting which have open *Issues* and open *Pull Requests* (PRs).

Use the `--filter` flag to only see repos with open Issues or PRs.

```
$ ghportfolio list

+---------------------------------------+-------------+----------+
|                 NAME                  | OPEN ISSUES | OPEN PRS |
+---------------------------------------+-------------+----------+
| the-general-problem-solver            |  !          |          |
| awesome-AI-academia                   |  !          |  !       |
| sudkamp-langs-machines-java           |  !          |          |
| sudkamp-langs-machines-python         |             |          |
| mAIcroft                              |             |          |
------------------------------------------------------------------
```

#### `interest`

Shows a timeline of when a particular repository was *Starred*, *Forked*, or *Cloned*. Use the `--chart` flag to get Sparkline visualisation.

```
$ ghportfolio interest google-rules-of-machine-learning

|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|
Stars:  |  |  |  |  |  | 1| 1|  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  | 1| 1|  |
Forks:  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |
Clones: |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |  |
```

```
$ ghportfolio interest --chart google-rules-of-machine-learning

|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|Sa|Fr|Th|We|Tu|Mo|Su|
Stars:  ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁██▁▁▁  min: 0 max: 1
Forks:  ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁  min: 0 max: 0
Clones: ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁  min: 0 max: 0
```

## Install

`go install github.com/thundergolfer/ghportfolio`

You can also download the relevant release for your OS from [`ghportfolio/releases`]( https://github.com/thundergolfer/ghportfolio/releases).

#### Uninstalling

You can just remove the `github.com/thundergolfer/ghportfolio` directory under `$GOPATH/bin/`

## Build

`go build` in root directory
