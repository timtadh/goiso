#include <stdlib.h>
#include <stdio.h>
#include <assert.h>
#include "graph.hh"
extern "C" {
#include "bliss_C.h"
}

/*
	Copyright (c) 2006-2011 Tommi Junttila
	Copyright (c) 2016 Tim Henderson
	Released under the GNU General Public License version 3.
	
	This file is part of bliss.
	
	bliss is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License version 3
	as published by the Free Software Foundation.
	
	bliss is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	
	You should have received a copy of the GNU General Public License
	along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
*/

extern "C"
struct bliss_graph_struct {
	bliss::Digraph* g;
};

extern "C"
struct bliss_stats_struct {
	/**
	 * An approximation (due to possible rounding errors) of
	 * the size of the automorphism group.
	 */
	long double group_size_approx;
	/** The number of nodes in the search tree. */
	long unsigned int nof_nodes;
	/** The number of leaf nodes in the search tree. */
	long unsigned int nof_leaf_nodes;
	/** The number of bad nodes in the search tree. */
	long unsigned int nof_bad_nodes;
	/** The number of canonical representative updates. */
	long unsigned int nof_canupdates;
	/** The number of generator permutations. */
	long unsigned int nof_generators;
	/** The maximal depth of the search tree. */
	unsigned long int max_level;
};

extern "C"
struct bliss_edge_struct {
	unsigned int Src;
	unsigned int Targ;
};

extern "C"
int
bliss_construct_and_canonize(unsigned int nodes[], int len_nodes, BlissEdge edges[], int len_edges, unsigned int perm[]) {
	BlissGraph *G;
	int i;
	const unsigned int * p;
	if (len_nodes <= 0 || len_edges < 0) {
		return 1;
	}
	if (nodes == NULL || perm == NULL) {
		return 2;
	}
	if (len_edges != 0 && edges == NULL) {
		return 3;
	}
	G = bliss_new(0);
	for (i = 0; i < len_nodes; i++) {
		bliss_add_vertex(G, nodes[i]);
	}
	for (i = 0; i < len_edges; i++) {
		bliss_add_edge(G, edges[i].Src, edges[i].Targ);
	}
	p = bliss_find_canonical_labeling(G, NULL, NULL, NULL);
	if (p == NULL) {
		return 4;
	}
	for (i = 0; i < len_nodes; i++) {
		perm[i] = p[i];
	}
	bliss_release(G);
	return 0;
}

extern "C"
BlissGraph *bliss_new(const unsigned int n)
{
	BlissGraph *graph = new bliss_graph_struct;
	assert(graph);
	graph->g = new bliss::Digraph(n);
	assert(graph->g);
	return graph;
}

extern "C"
BlissGraph *bliss_read_dimacs(FILE *fp)
{
	bliss::Digraph *g = bliss::Digraph::read_dimacs(fp);
	if(!g) {
		return 0;
	}
	BlissGraph *graph = new bliss_graph_struct;
	assert(graph);
	graph->g = g;
	return graph;
}

extern "C"
void bliss_write_dimacs(BlissGraph *graph, FILE *fp)
{
	assert(graph);
	assert(graph->g);
	graph->g->write_dimacs(fp);
}

extern "C"
void bliss_release(BlissGraph *graph)
{
	assert(graph);
	assert(graph->g);
	delete graph->g; graph->g = 0;
	delete graph;
}

extern "C"
void bliss_write_dot(BlissGraph *graph, FILE *fp)
{
	assert(graph);
	assert(graph->g);
	graph->g->write_dot(fp);
}

extern "C"
unsigned int bliss_get_nof_vertices(BlissGraph *graph)
{
	assert(graph);
	assert(graph->g);
	return graph->g->get_nof_vertices();
}

extern "C"
unsigned int bliss_add_vertex(BlissGraph *graph, unsigned int l)
{
	assert(graph);
	assert(graph->g);
	return graph->g->add_vertex(l);
}

extern "C"
void bliss_add_edge(BlissGraph *graph, unsigned int v1, unsigned int v2)
{
	assert(graph);
	assert(graph->g);
	graph->g->add_edge(v1, v2);
}

extern "C"
int bliss_cmp(BlissGraph *graph1, BlissGraph *graph2)
{
	assert(graph1);
	assert(graph1->g);
	assert(graph2);
	assert(graph2->g);
	return (*graph1->g).cmp(*graph2->g);
}

extern "C"
unsigned int bliss_hash(BlissGraph *graph)
{
	assert(graph);
	assert(graph->g);
	return graph->g->get_hash();
}

extern "C"
BlissGraph *bliss_permute(BlissGraph *graph, const unsigned int *perm)
{
	assert(graph);
	assert(graph->g);
	assert(graph->g->get_nof_vertices() == 0 || perm);
	BlissGraph *permuted_graph = new bliss_graph_struct;
	assert(permuted_graph);
	permuted_graph->g = graph->g->permute(perm);
	return permuted_graph;
}

extern "C"
void
bliss_find_automorphisms(BlissGraph *graph,
			 void (*hook)(void *user_param,
				      unsigned int n,
				      const unsigned int *aut),
			 void *hook_user_param,
			 BlissStats *stats)
{
	bliss::Stats s;
	assert(graph);
	assert(graph->g);
	graph->g->find_automorphisms(s, hook, hook_user_param);

	if(stats)
	{
		stats->group_size_approx = s.get_group_size_approx();
		stats->nof_nodes = s.get_nof_nodes();
		stats->nof_leaf_nodes = s.get_nof_leaf_nodes();
		stats->nof_bad_nodes = s.get_nof_bad_nodes();
		stats->nof_canupdates = s.get_nof_canupdates();
		stats->nof_generators = s.get_nof_generators();
		stats->max_level = s.get_max_level();
	}
}


extern "C"
const unsigned int *
bliss_find_canonical_labeling(BlissGraph *graph,
			      void (*hook)(void *user_param,
					   unsigned int n,
					   const unsigned int *aut),
			      void *hook_user_param,
			      BlissStats *stats)
{
	bliss::Stats s;
	const unsigned int *canonical_labeling = 0;
	assert(graph);
	assert(graph->g);

	canonical_labeling = graph->g->canonical_form(s, hook, hook_user_param);

	if(stats)
	{
		stats->group_size_approx = s.get_group_size_approx();
		stats->nof_nodes = s.get_nof_nodes();
		stats->nof_leaf_nodes = s.get_nof_leaf_nodes();
		stats->nof_bad_nodes = s.get_nof_bad_nodes();
		stats->nof_canupdates = s.get_nof_canupdates();
		stats->nof_generators = s.get_nof_generators();
		stats->max_level = s.get_max_level();
	}

	return canonical_labeling;
}
