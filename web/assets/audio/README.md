# Animal recordings

`species/<genus>-<species>.mp3` is a real recording of that species, trimmed to a
short clip. Each file is listed with its source, author and licence in
`web/views/components/sounds.go`, which also feeds the /credits page. Sources are
Wikimedia Commons and NOAA Fisheries (public domain).

To add one, drop the mp3 in `species/` and add an entry to `recordings`.
