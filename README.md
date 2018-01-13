# ghportfolio [![BuildStatus](https://travis-ci.com/thundergolfer/ghportfolio.svg?token=yHGWQ42iK2BPk1FjaUMc&branch=master)](https://travis-ci.com/thundergolfer/ghportfolio)
A minimal CLI in Golang for catching up on the activity and health of your public Github projects

## Why Use It?

After switching to a team that used the Bitbucket ecosystem, I lost touch with the state of my open-source projects as I was no longer arriving at the Github homepage everyday.

I made this small CLI tool to quickly check-up on projects, seeing whether had received any recent interest from other users, or whether there was some issue or code contribution that had popped up on a project without me noticing.

Hope you may also find it useful :)

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

You can also use the `--totals` flag to get a cumulative total time-series of repo stars.

```
$ ./ghportfolio interest interview-with-python --totals

92.4│
    │                                                                                            •••
    │                                                                        •••••••••••••••••••••
    │                                                             ••••••••••••
    │                                                       ••••••
    │                                          ••••••••••••••
    │                                     •••••
    S                      ••••••••••••••••
    t     ••••••••••••••••••
    a    ••
    r    •
    s •••
    │••
    │•
    │•
    │•
    │•
    │•
0   │-------------------------------------------Time------------------------------------------------
     0.0                                                                                       484.0
 current: 84
 ```

## Install

`go install github.com/thundergolfer/ghportfolio`

You can also download the relevant release for your OS from [`ghportfolio/releases`]( https://github.com/thundergolfer/ghportfolio/releases).

#### Uninstalling

You can just remove the `github.com/thundergolfer/ghportfolio` directory under `$GOPATH/bin/`

## Build

`go build` in root directory
