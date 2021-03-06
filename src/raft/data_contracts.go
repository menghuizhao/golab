package raft

import (
	"labrpc"
	. "menghuibasic"
	"sync"
	"time"
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// Persistent state on all servers
	currentTerm int
	votedFor    NullInt
	log         []LogEntry

	// Volatile state on all servers
	state       State
	commitIndex int
	lastApplied int

	// Volatile state on leaders
	nextIndex  []int
	matchIndex []int

	// Utility
	electionTimer   *time.Timer
	electionTimeout time.Duration
	roleChan        chan State
	myVotes         int
}

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	// field names must start with capital letters!
	CandidateTerm int // Candidate's term
	CandidateId   int // Candidate requesting vote
	LastLogIndex  int // Index of candidate's last log entry
	LastLogTerm   int // Term of candidate's last log entry
}

type RequestVoteReply struct {
	// Your data here (2A).
	// field names must start with capital letters!
	Term        int  // Current Term for candidate to update itself
	VoteGranted bool // True if candidates is granted this vote
}

type AppendEntriesArgs struct {
	Term         int        //leader’s term
	LeaderId     int        //so follower can redirect clients
	PrevLogIndex int        //index of log entry immediately preceding new ones
	PrevLogTerm  int        //term of prevLogIndex entry
	Entries      []LogEntry //log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        //leader’s commitIndex
}

type AppendEntriesReply struct {
	Term    int  //currentTerm, for leader to update itself
	Success bool //true if follower contained entry matching prevLogIndex and prevLogTerm
}

type LogEntry struct {
	Command string
	Term    int
}

/*
	Constants
*/
type State int

const (
	Undefined State = iota // value --> 0
	Follower               // value --> 1
	Candidate              // value --> 2
	Leader                 // value --> 3
)

func (state State) String() string {
	switch state {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

const heartbeatInterval = 101 * time.Millisecond
