package mechlibrary

import (
    "fmt"
    "math/rand"

    "appengine"
    "appengine/datastore"
    "appengine/memcache"
)

type voteConfig struct {
    Shards int
}

type voteShard struct {
    Name  string
    VoteCount int
}

const (
    defaultShards = 20
    configKind    = "MechanicVoteShardConfig"
    shardKind     = "MechanicVoteShard"
)

func memcacheKey(name string) string {
    return shardKind + ":" + name
}

// Count retrieves the value of the named counter.
func Count(c appengine.Context, name string) (int, bool, error) {
    //c.Errorf("counting: %v", name)
    total := 0
    mkey := memcacheKey(name)
    if _, err := memcache.JSON.Get(c, mkey, &total); err == nil {
        //c.Errorf("returning cached: %v", total)
        return total, false, nil
    }
    q := datastore.NewQuery(shardKind).Filter("Name =", name)
    for t := q.Run(c); ; {
        var s voteShard
        _, err := t.Next(&s)
        if err == datastore.Done {
            break
        }
        if err != nil {
            return total, true, err
        }
        total += s.VoteCount
    }
    memcache.JSON.Set(c, &memcache.Item{
        Key:        mkey,
        Object:     &total,
        Expiration: 60,
    })
    //c.Errorf("returning not cached: %v", total)
    return total, true, nil
}

// Increment increments the named counter.
func Increment(c appengine.Context, name string) error {
    //c.Errorf("incrementing: %v", name)
    // Get counter config.
    var cfg voteConfig
    ckey := datastore.NewKey(c, configKind, name, 0, nil)
    err := datastore.RunInTransaction(c, func(c appengine.Context) error {
        err := datastore.Get(c, ckey, &cfg)
        if err == datastore.ErrNoSuchEntity {
            cfg.Shards = defaultShards
            _, err = datastore.Put(c, ckey, &cfg)
        }
        return err
    }, nil)
    if err != nil {
        return err
    }
    var s voteShard
    err = datastore.RunInTransaction(c, func(c appengine.Context) error {
        shardName := fmt.Sprintf("%s-shard%d", name, rand.Intn(cfg.Shards))
        key := datastore.NewKey(c, shardKind, shardName, 0, nil)
        err := datastore.Get(c, key, &s)
        // A missing entity and a present entity will both work.
        if err != nil && err != datastore.ErrNoSuchEntity {
            return err
        }
        s.Name = name
        s.VoteCount++
        _, err = datastore.Put(c, key, &s)
        return err
    }, nil)
    if err != nil {
        return err
    }
    memcache.IncrementExisting(c, memcacheKey(name), 1)
    return nil
}
