// Computes canonical permutations (labeling) of directed labeled graphs. Useful
// for comparing graphs and answering graph isomorphism questions. If you need
// to know whether two graphs are equivalent this is the library for you!
//
// A wrapper around [bliss](http://www.tcs.hut.fi/Software/bliss/) for graph
// isomorphism testing and canonical labeling. Bliss is the work Tommi Junttila and
// Petteri Kaski. You should cite their papers:
//
// - Tommi Junttila and Petteri Kaski. Engineering an efficient canonical labeling
//   tool for large and sparse graphs. In: Proceedings of the Ninth Workshop on
//   Algorithm Engineering and Experiments (ALENEX07), pages 135-149, SIAM, 2007.
//
// - Tommi Junttila and Petteri Kaski. Conflict Propagation and Component Recursion
//   for Canonical Labeling. In: Proceedings of the 1st International ICST
//   Conference on Theory and Practice of Algorithms (TAPAS 2011), Springer, 2011
//
// I have made a few modifications to their library mainly in the C interface side
// so I can interface more easily with Go.
//
//
// All work is licensed under the GPL Version 3. Please respect the license.
//
// Copyright Holders:
//
//       Copyright (c) 2006-2011 Tommi Junttila (C++ Bliss Code)
//       Copyright (c) 2014-2016 Tim Henderson (Go Wrapper)
package bliss
