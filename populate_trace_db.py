#!/usr/bin/env python

# A utility that reads in local trace CSVs and populate the trace data
# into Redis.
# Expecting the CSV format to be the same specified in receiver.cpp.

import sys
import redis

class Trace:
  def __init__(self, line):
    fields = line.split(",")
    self.id = fields[0]
    self.name = fields[1]
    self.span_id = fields[2]
    self.span_parent_id = fields[3]
    self.from_upid = fields[4]
    self.to_upid = fields[5]
    self.time_nanos = int(fields[6])
    self.stage = fields[7]

r = redis.StrictRedis(host='localhost', port=6379, db=0)

while 1:
  line = sys.stdin.readline()
  if not line: break
  trace = Trace(line)
  if not r.hexists('traces', trace.id):
    r.hset('traces', trace.id, trace.time_nanos)
  else:
    time = r.hget('traces', trace.id)
    if time > trace.time_nanos:
      r.hset('traces', trace.id, trace.time_nanos)

  i = r.incr(trace.id + "-count")

  list_key = trace.id + "-" + str(i)
  r.hset(list_key, 'name', trace.name)
  r.hset(list_key, 'span_id', trace.span_id)
  r.hset(list_key, 'span_parent_id', trace.span_parent_id)
  r.hset(list_key, 'from_upid', trace.from_upid)
  r.hset(list_key, 'to_upid', trace.to_upid)
  r.hset(list_key, 'time_nanos', trace.time_nanos)
  r.hset(list_key, 'stage', trace.stage)




