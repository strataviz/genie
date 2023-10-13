package resources

import (
	"ctx.sh/strata"
	"github.com/go-logr/logr"
	"stvz.io/genie/pkg/resources/integer_range"
	"stvz.io/genie/pkg/resources/ipaddr"
	"stvz.io/genie/pkg/resources/list"
	"stvz.io/genie/pkg/resources/random_string"
	"stvz.io/genie/pkg/resources/timestamp"
	"stvz.io/genie/pkg/resources/uuid"
)

type Resource interface {
	Get() string
}

type Options struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type Resources struct {
	Lists         map[string]Resource
	IntegerRanges map[string]Resource
	RandomStrings map[string]Resource
	Uuids         map[string]Resource
	Timestamps    map[string]Resource
	Maps          map[string]Resource
	IPAddrs       map[string]Resource
}

func Parse(block Config, opts *Options) (*Resources, error) {
	integerRanges, err := parseIntegerRanges(block)
	if err != nil {
		return nil, err
	}

	lists, err := parseLists(block)
	if err != nil {
		return nil, err
	}

	randomStrings, err := parseRandomStrings(block)
	if err != nil {
		return nil, err
	}

	timestamps, err := parseTimestamps(block)
	if err != nil {
		return nil, err
	}

	uuids, err := parseUuids(block)
	if err != nil {
		return nil, err
	}

	ipaddrs, err := parseIPAddrs(block)
	if err != nil {
		return nil, err
	}

	return &Resources{
		IntegerRanges: integerRanges,
		RandomStrings: randomStrings,
		Lists:         lists,
		Timestamps:    timestamps,
		Uuids:         uuids,
		IPAddrs:       ipaddrs,
	}, nil
}

func parseIntegerRanges(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.IntegerRanges {
		out[k] = integer_range.New(v)
	}

	return out, nil
}

func parseRandomStrings(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.RandomStrings {
		out[k] = random_string.New(v)
	}

	return out, nil
}

func parseLists(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.Lists {
		out[k] = list.New(v)
	}

	return out, nil
}

func parseTimestamps(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.Timestamps {
		out[k] = timestamp.New(v)
	}

	return out, nil
}

func parseUuids(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.Uuids {
		out[k] = uuid.New(v)
	}

	return out, nil
}

func parseIPAddrs(res Config) (map[string]Resource, error) {
	out := make(map[string]Resource)

	for k, v := range res.IPAddrs {
		out[k] = ipaddr.New(v)
	}

	return out, nil
}

func (r *Resources) Get(rtype string, name string) (Resource, error) {
	switch rtype {
	case "list":
		return r.GetList(name)
	case "integer_range":
		return r.GetIntegerRange(name)
	case "random_string":
		return r.GetRandomString(name)
	case "uuid":
		return r.GetUuid(name)
	case "timestamp":
		return r.GetTimestamp(name)
	case "map":
		return r.GetMap(name)
	case "ipaddr":
		return r.GetIPAddr(name)
	default:
		return nil, InvalidResourceTypeError
	}
}

func (r *Resources) MustGet(rtype string, name string) Resource {
	res, err := r.Get(rtype, name)
	if err != nil {
		panic(err)
	}

	return res
}

func (r *Resources) GetList(name string) (Resource, error) {
	if resource, ok := r.Lists[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetIntegerRange(name string) (Resource, error) {
	if resource, ok := r.IntegerRanges[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetRandomString(name string) (Resource, error) {
	if resource, ok := r.RandomStrings[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetUuid(name string) (Resource, error) {
	if resource, ok := r.Uuids[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetTimestamp(name string) (Resource, error) {
	if resource, ok := r.Timestamps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetMap(name string) (Resource, error) {
	if resource, ok := r.Maps[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}

func (r *Resources) GetIPAddr(name string) (Resource, error) {
	if resource, ok := r.IPAddrs[name]; ok {
		return resource, nil
	}

	return nil, NotFoundError
}
