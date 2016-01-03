package main

import (
	decoder "github.com/minchao/googlemaps-share-link-decoder"
)

type ServiceMiddleware func(decoder.Service) decoder.Service
