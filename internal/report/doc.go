// Package report renders an Ekman spiral computation as human-readable text
// for the command line.
//
// The output has three parts:
//
//  1. a header restating every input the case resolved to (wind stress,
//     density, Coriolis parameter, eddy viscosity, optional finite depth);
//  2. the headline results — surface current magnitude and heading, Ekman
//     depth De, and the depth-integrated transport magnitude and heading;
//  3. a table of the sampled spiral, one row per depth, with east/north
//     velocity components, speed, heading and veer relative to the surface
//     current.
//
// All angles are printed in degrees clockwise from true north, matching the
// physics layer's convention. The package owns no computation and performs no
// validation; it only formats a result that internal/spiral already produced.
package report
