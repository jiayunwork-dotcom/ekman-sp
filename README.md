# ekman-sp

A small, self-contained command-line calculator for the **steady-state ocean
surface Ekman spiral**. Given a wind stress, seawater density, the Coriolis
parameter (or a latitude) and an eddy viscosity, it computes the surface
current magnitude and direction, the Ekman depth, the depth-integrated Ekman
transport, and the turning velocity profile with depth.

This is the wind-driven frictional kernel of the surface ocean — not a sea
state dashboard and not tidal or wave analysis (tides and linear waves live in
separate tools).

## What it computes

With the vertical coordinate positive upward from the sea surface and the
Coriolis parameter `f` positive in the northern hemisphere, the steady
f-plane equations are:

```
K d²u/dz² = -f v
K d²v/dz² =  f u
```

The deep-water solution is a damped, turning spiral. The Ekman depth scale is

```
delta = sqrt(2K/|f|)         (m)
De    = pi * delta           (Ekman depth, m)
```

the surface current magnitude is

```
|V0| = tau / (rho * sqrt(K*|f|))     (m/s)
```

and the depth-integrated volume transport is

```
Me = tau / (rho * |f|)               (m²/s)
```

In the northern hemisphere the surface current is deflected **45° to the
right of the wind** and turns clockwise with depth; in the southern hemisphere
both deflections reverse. The transport lies 90° to the right of the wind in
the north and 90° to the left in the south.

## Inputs, outputs, and boundaries

- `tau` — wind stress magnitude (N/m²), required and must be positive.
- `wind_dir_deg` — direction the wind blows toward, degrees clockwise from
  north (default: 90, an easterly wind).
- `rho` — seawater density (kg/m³, default 1025).
- `K` — vertical eddy viscosity (m²/s, default 0.05).
- `f` **or** `latitude_deg` — the Coriolis parameter (1/s) or the latitude in
  degrees, from which `f = 2·Ω·sin(lat)` is derived. At least one is
  required; supplying both is allowed only when they agree. Latitude 0 or
  `f = 0` (the equator) is an error.
- `depth_m` — optional finite water depth (m). A depth smaller than the Ekman
  depth `De` is rejected: the classic infinite-depth spiral no longer applies
  in such a thin layer. Deeper cases are computed with the infinite-depth
  solution and the residual boundary deviation is declared in the output.
- `npoints` — profile sample count (default 21).

Invalid input produces a clear error on standard error and a non-zero exit
code; nothing is silently ignored.

## Usage

```bash
go run . profile example/midlat-wind.json
```

`example/midlat-wind.json` is a mid-latitude westerly wind case at 45°N. The
report prints the surface current magnitude and heading (135°, i.e. 45° right
of the easterly wind), the Ekman depth `De`, the transport magnitude and
heading (180°, 90° right of the wind), and the depth profile table:

```text
Ekman spiral — steady-state wind-driven surface layer
  inputs
    wind stress tau        : 0.200 N/m^2 toward 90.0 deg (E)
    density rho            : 1025.0 kg/m^3
    Coriolis f             : +1.031e-04 1/s (lat 45.0 deg, northern hemisphere)
    eddy viscosity K       : 0.050 m^2/s
    water depth            : infinite (classic Ekman solution)
    profile samples        : 21
  results
    surface speed |V0|    : 0.0859 m/s  (= tau/(rho*sqrt(K*|f|)))
    surface heading       : 135.0 deg (SE)  (45.0 deg right of wind)
    surface velocity (u,v): (0.0608, -0.0608) m/s
    Ekman scale delta     : 31.14 m  (= sqrt(2K/|f|))
    Ekman depth De        : 97.83 m  (= pi*sqrt(2K/|f|))
    transport |Me|        : 1.8921 m^2/s  (= tau/(rho*|f|))
    transport heading     : 180.0 deg (S)  (90.0 deg right of wind)
    transport vector (u,v): (0.0000, -1.8921) m^2/s
    wind direction        : 90.0 deg (E) (for reference)
  depth profile
  depth (m) u (m/s)  v (m/s)  |V| (m/s) heading (deg)veer (deg)
  ...
```

Other cases:

```bash
go run . profile example/south-hemi-wind.json   # southern hemisphere: leftward deflection
go run . profile example/north-trade.json       # tropical easterlies at 15°N
```

## Cross-rule checks

The implementation obeys the scaling relations that tie the spiral together:

- doubling `tau` doubles both `|V0|` and `Me`;
- doubling `|f|` divides `De` by `√2` and reduces `|V0|`;
- doubling `K` multiplies `De` by `√2`;
- reversing the wind reverses the entire spiral;
- flipping the sign of `f` swaps the left/right deflections.

## Build & test

```bash
go build ./...
go test ./...
```

## License

MIT — see [LICENSE](LICENSE).
