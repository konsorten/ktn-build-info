package ver

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/consul/api"
)

func RetrieveBuildFromConsul(consulUrl string, kvProjectRoot string, vi *VersionInformation) error {
	// check version info
	if kvProjectRoot == "" {
		return fmt.Errorf("Consul KV project root path may not be empty")
	}

	if vi.Revision == "" {
		return fmt.Errorf("No revision information available")
	}

	// parse the url
	cu, err := url.Parse(consulUrl)

	if err != nil {
		return fmt.Errorf("Failed to parse consul URL: %v", err)
	}

	// assign config
	cfg := api.DefaultConfig()

	cfg.Address = cu.Host
	cfg.Scheme = cu.Scheme
	cfg.Datacenter = cu.Path[1:] // trim first slash

	if cu.User != nil {
		if t, ok := cu.User.Password(); ok {
			cfg.Token = t
		}
	}

	// connect
	client, err := api.NewClient(cfg)

	if err != nil {
		return fmt.Errorf("Failed to connect to consul: %v", err)
	}

	kv := client.KV()

	// create session
	session, _, err := client.Session().Create(&api.SessionEntry{TTL: "60s", Behavior: "delete"}, nil)

	if err != nil {
		return fmt.Errorf("Failed to create session: %v", err)
	}

	defer client.Session().Destroy(session, nil)

	// acquire revision
	revPath := fmt.Sprintf("%v/revs/%v", kvProjectRoot, vi.Revision)
	lockPath := fmt.Sprintf("%v/~lock", revPath)

	locked, _, err := kv.Acquire(&api.KVPair{Key: lockPath, Session: session}, nil)

	if err != nil {
		return fmt.Errorf("Failed to acquire lock: %v", err)
	}

	if !locked {
		return fmt.Errorf("Failed to acquire lock on KV path: %v", lockPath)
	}

	defer kv.Release(&api.KVPair{Key: lockPath, Session: session}, nil)

	// retrieve any existing
	build, _, err := kv.Get(fmt.Sprintf("%v/build", revPath), nil)

	if err != nil {
		return fmt.Errorf("Failed to retrieve existing build id: %v", err)
	}

	if build != nil {
		b, err := strconv.Atoi(string(build.Value))

		if err != nil {
			return fmt.Errorf("Failed to parse existing build id: %v", err)
		}

		vi.Build = b

		return nil
	}

	// get new build id
	for true {
		nextIdPath := fmt.Sprintf("%v/nextId", kvProjectRoot)

		nextId, _, err := kv.Get(nextIdPath, nil)

		if err != nil {
			return fmt.Errorf("Failed to get next build id: %v", err)
		}

		if nextId != nil {
			nextId.Session = session

			// parse the id
			vi.Build, err = strconv.Atoi(string(nextId.Value))

			if err != nil {
				return fmt.Errorf("Failed to parse next build id: %v", err)

			}

			// increment the id
			nextId.Value = []byte(strconv.Itoa(vi.Build + 1))
		} else {
			// use 1 as first ID, and register 2 as the next one
			vi.Build = 1

			nextId = &api.KVPair{Key: nextIdPath, Value: []byte("2"), Session: session}
		}

		// update the id
		ok, _, err := kv.CAS(nextId, nil)

		if ok {
			break
		}
	}

	// write the revision
	_, err = kv.Put(&api.KVPair{
		Key:     fmt.Sprintf("%v/build", revPath),
		Value:   []byte(strconv.Itoa(vi.Build)),
		Session: session,
	}, nil)

	if err != nil {
		return fmt.Errorf("Failed to write build revision: %v", err)
	}

	// done
	return nil
}
