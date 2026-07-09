<p align="center">
  <img src="assets/lucario_banner.jpg" alt="banner" width="400"/>
</p>

# Lucario

Lucario is a physical Write-Ahead Logging (WAL) implementation designed for custom database storage engines. It provides durability and crash recovery by recording page-level modifications before they are applied to disk and replaying committed operations after a crash.

## Motivation

Database operations such as inserting a key-value pair appear atomic from an application's perspective. Internally, however, a single operation may:

* Modify multiple pages
* Split B+ Tree nodes
* Allocate new pages
* Update metadata pages

A crash occurring in the middle of these modifications can leave the database in an inconsistent state. For example, a leaf node may have already been split while its parent has not yet been updated, effectively orphaning an entire subtree.

Lucario solves this problem by recording modifications in a Write-Ahead Log and replaying only complete operations during recovery.


## Features

* Physical page-level redo logging
* Log Sequence Number (LSN) tracking
* [Iterator-based WAL record replay](https://github.com/Adarsh-Kmt/Lucario/blob/master/wal_iterator.go) during crash recovery
* Atomic recovery using `BeginOperation` and `CommitOperation` records
* Generic log record serialization framework


## Architecture

Every logical storage engine operation is represented by a sequence of WAL records:

```text
┌──────────────────────┐
│   BeginOperation     │
├──────────────────────┤
│   WAL Record 1       │
│   WAL Record 2       │
│         .            │
│         .            │
│   WAL Record N       │
├──────────────────────┤
│   CommitOperation    │
└──────────────────────┘
```

The storage engine follows the Write-Ahead Logging protocol:

```text
Generate WAL Records
          │
          ▼
Append WAL Records
          │
          ▼
Flush WAL to Disk
          │
          ▼
Apply Page Modifications
```

This guarantees that recovery always has enough information to reconstruct the database state after a crash.


## WAL Record Format

Each WAL record consists of:

| Field          | Description                                  |
| -------------- | -------------------------------------------- |
| WAL Record length | Length of the WAL Record                  |
| LSN            | Monotonically increasing Log Sequence Number |
| Operation      | Type of operation being logged               |
| Payload Length | Size of serialized payload                   |
| Payload        | Operation-specific data                      |
| WAL Record length | Length of the WAL Record                  |


## Recovery Algorithm

Recovery is implemented using [iterator-based replay](https://github.com/Adarsh-Kmt/Lucario/blob/master/wal_iterator.go).

### Recovery Steps

1. Sequentially scan the WAL.
2. Buffer records between `BeginOperation` and `CommitOperation`.
3. If a `CommitOperation` is encountered:

   * Replay all buffered records.
4. If the end of the WAL is reached before a `CommitOperation`:

   * Discard the buffered records.

This ensures that partially completed operations are never replayed.


## Idempotent Recovery

Each database page maintains a `pageLSN`, representing the most recent WAL record applied to that page.

During recovery:

```go
if pageLSN < recordLSN {
    apply(record)
}
```

This makes recovery idempotent and allows the recovery process to be safely restarted multiple times.


## Example

A B+ Tree insert that causes cascading splits may generate the following WAL sequence:

```text
BeginOperation
├── SplitLeafNode
├── InsertInternalNodeEntry
├── SplitInternalNode
├── InsertInternalNodeEntry
└── UpdateRootNodePageId
CommitOperation
```

If the database crashes before `CommitOperation` is written:

```text
BeginOperation
├── SplitLeafNode
└── InsertInternalNodeEntry
CRASH
```

the buffered records are discarded during recovery, preventing the B+ Tree from entering an inconsistent intermediate state.


## Recovery Flow

```text
Database Startup
        │
        ▼
Open WAL
        │
        ▼
Sequential WAL Scan
        │
        ▼
Collect Records Between
BeginOperation and CommitOperation
        │
        ▼
Replay Committed Operations
        │
        ▼
Flush Recovered Pages
        │
        ▼
Database Ready
```

## Future Work

* Physical replication using streamed WAL records



### Resources I Used
- [Kevin Sookocheff's Blog](https://sookocheff.com/post/databases/write-ahead-logging/)
- [Avtandil Ushikishvili's Blog](https://avtandili.ge/posts/trying-to-understand-wal-under-the-hood/)
