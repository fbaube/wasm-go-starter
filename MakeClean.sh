rm ~/.wash/dev/*
wash drain dev
rm ~/.wash/package_cache/*
killall wasm nats-server wasmcloud_host wadm wash
