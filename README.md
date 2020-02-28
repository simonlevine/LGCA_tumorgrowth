# LGCA_tumorgrowth
Tumor Metastasis Simulation via Lattice-Gas Cellular Automata
 
Abstract
Cancer remains a devastating disease even in the modern era of medicine. While the molecular characteristics of cancer growth have been studied for some time, it is only relatively recently that methods within computational biology have been applied to this aspect of oncology. As computational capabilities progress, tumor growth may prove to be pragmatically studied in silico.
One such method study leverages the conceptual foundation of the cellular automata model and extends it using physical concepts, including the Cellular Potts model and gaseous fluid theory as applied to Euclidean lattices of cells. 

Purpose: Implementation of a cellular automata model that is modular with respect to neighborhoods, dimensions, plotting, and physical constants.

§1. Introduction
This cellular automata simulation leverages existing models first developed to simulate gaseous particles and applies the same computational logic to growing tumor cells. The so-called Lattice Gas Cellular Automata (LGCA) uses a finite Euclidean lattice of cell “sites”, such that each site has interactivity with neighboring sites via predefined neighborhood bounds (Von Neumann; this particular choice is further explained). Additionally, proxies for cell-to-cell interaction within neighborhoods are modeled via like-type proximity. This approach was first implemented computationally and plotted in two dimensions, and then implemented in three dimensions.

§2. Methods

§2.1 Two-Dimensional LGCA
Lattice-Gas Cellular Automata (in two dimensions) leverages existing principles of cellular automata  and adds movement of particles – cells with states this case – through discrete velocity vectors on a two-dimensional matrix. Depending on its local Von Neumann neighborhood, a given cell’s likely state in the next timestep of the simulation is computed (the “reactive step”.). Then, a cell’s propensity to move of propagate is applied via stochastic movement along discrete vectors to neighboring sites (Von Neumann neighborhoods, whereby each site on the lattice is associated with five possible channels of movement).
To this end, we wrote a Go program that declares a two-dimensional lattice of cell structs, whereby each cell bears a state in the form of a string, a location in the form of an ordered pair of integers, a velocity direction (for movement and propagation purposes), and an array of pointers to cells of its local neighborhood.

§2.2 Reactive Step
Cells may transition from state to state  . That is, we define “cells” within the lattice as either cancerous, healthy, or necrotic. Cancerous cells can either be in a state of quiescence (non-propagative, but not necrotic either) or in a proliferative state (whereby cell division will take place and a new cancer cell will propagate at the next timestep). Necrotic cells are defined as previously cancerous cells that died due to lack of available resources (space, in this case). Healthy cells are simulated as sites on the lattice not occupied by other cells, or sites upon which other cells can invade and spread. Apoptotic cells were additionally modeled though not implemented in the runnable simulation for simplicity, though modular, workable code is included as “commented” code.
Cell state transitions are computed probabilistically via a simplified adaptation of Lattice-Boltzmann Energy theory. Three cell coupling coefficient  parameters are employed to this effect: K_cc, K_nn, and K_cn. These represent modeling constants proportional to the strength of membrane coupling between cancer cells (quiescent included), necrotic-necrotic cell interaction, and cancerous-necrotic interaction.

Next, these constants are imputed into Lattice-Boltzmann energy factor equations defined for cells of type proliferative, quiescent, and necrotic, together with the count of cells of each type in the current neighborhood (C, N, for cancerous and necrotic, respectively):

	E_p=-[0.50(C(C+1)) K_cc+N(N-1) K_nn )+(C+1)N K_cn]
	E_q=-[0.50((C-1)(C)) K_cc+N(N-1) K_nn )+(C)N K_cn]
	E_n=-[0.50((C-1)(C-2)) K_cc+N(N+1) K_nn )+(C-1)(N+1) K_cn]

In other words, one of the following reactions is possible for between any two given timesteps at site i,j of the lattice  :

	Proliferation: C_(i,j)→C_(i,j)+1 ; N_(i,j)→N_(i,j) .
	Quiescence: C_(i,j)→C_(i,j) ; N_(i,j)→N_(i,j) (no change).
	Necrosis: C_(i,j)→C_(i,j)-1 ; N_(i,j)→N_(i,j)+1 .

Note that these Boltzmann models are highly related, proportional abstractions of the Hamiltonian energy model employed in the generalized Cellular Potts  automata model .
Next, these computed Boltzmann energies are imputed into a simplified ,  lattice-gas probability model, such that proxies for proliferative, quiescence, and necrosis likelihoods are obtained by comparing ratios of Boltzmann factors :

	P_proliferative=e^(E_p )/(e^(E_p )+e^(E_q )+e^(E_n ) )
	P_quiescent=e^(E_q )/(e^(E_p )+e^(E_q )+e^(E_n ) )
	P_necrosis=e^(E_n )/(e^(E_p )+e^(E_q )+e^(E_n ) )

States of each cell on the lattice are then updated according to these probabilities at each timestep.

§2.3 Movement Step
Following the computational “reactive step” above, each cell is primed for potential movement as a function of their probabilistic states. Propagative cancerous cells will divide, with nascent cells invading the adjacent neighborhood least dense in cancerous cells (modeling invasion), with parental cells remaining in a current site. Modeling principles of chemotaxis , necrotic cells will move toward regions most dense in other necrotic cells (the simulation plots a path of movement). Finally, quiescent cancer cells will not move. If a tie is obtained between cell-type densities in adjacent neighborhoods (e.g., if two or more adjacent neighborhoods contain the name number of C or N cells), one of these equivalent neighborhoods is chosen at random for cell movement and/or propagation.

Cells of relevant types are then moved synchronously (i.e., a single timestep change results in all cells updated across the lattice and within neighborhoods) at the given timestep, resulting in an updated lattice. Note that a Von Neumann cellular automata neighborhood (see Figure 1.) was chosen in this instance due to practicality of implementation  as well as reasons of physical verisimilitude .
 
Figure 1. 2-D Von Neumann neighborhood with adjacent lattice sites. Velocity vectors superimposed.

§2.4 Plotting

§2.5 Beyond 2 Dimensions
The implementation of this simulation is highly modular, and the code easily lends itself to expansion. One such expansion that was implemented was an extension to three-dimensional space. That is, a three-dimensional lattice of cells were defined on a row, column, and “aisle” basis. All computational and logistical two-dimensional functions were then altered to accommodate the new spatial arrangement of the simulation.

§2.6 Three-Dimensional Reactive Step
This step is logically equivalent to that of the two-dimensional implementation (see §2.2), only that cell states at a given timestep are defined by a three-dimensional Von Neumann neighborhood (see Figure 2.)
 
Figure 2. The 3-D Von Neumann arrangement, with velocity vectors of the center cell superimposed.

§2.7 Propagation Step in Three Dimensions
Cell movement and propagation is again equivalent logically to the two-dimensional implementation (see §2.3), only that the implementation required an extension into the “aisle” space, accommodating the extra dimension. Cell movement has two extra potential axes (see Figure 2.) for two additional discrete velocity vectors.

§2.8 Plotting in Three Dimensions
 
§3. Results
§3.1 Two-Dimensional Plot
 

Figure 3. Initial two-dimensional lattice seeded with a cancerous central neighborhood and adjacent sites. Units in cells (pixels).
 
Figure 4. Lattice after 50 generations of growth.

§3.2 Three-Dimensional Plot

 

Figure 5. Three-dimensional lattice similarly seeded with cancerous cells.


 

Figure 6. Three-dimensional lattice after 35 generations of the simulation (only cancerous cells plotted for simplicity).

§3.3 Discussion
The model was simulated on a two-dimensional lattice for 50 generations (see Figures 3,4.), and a three-dimensional model was simulated for 35 generations of growth (see Figures 5,6.). Additionally, for this simulation, coupling constant parameters were set to K_cc=K_nn=3.0; K_nc=1.0 per the advice of previous two-dimensional investigation . 
Though the rules employed in this automata are somewhat simple compared to true cellular dynamics, nonetheless tumor “cells” are modeled somewhat realistically. Competition for resources (spatial constraints) here between healthy (background) cells and proliferative cancer cells is apparent, as is the complex interplay between cells of difference states.

 
§4. Conclusions
This implementation represents a successful, novel implementation of lattice-gas cellular automata as applied to principles of tumor growth. Furthermore, this implementation can be run in virtually infinitely many ways simply by altering existing constraints. Additionally, the model is sufficiently modular to allow for simple expansion, such as the addition of apoptotic cells, healthy cell behavior, additional coupling constraints, and more complex neighborhood geometry.
Some limitations of this model are apparent, however. The implementation is compute-intensive and memory-heavy; more complex investigations would need to better utilize principles of parallelism to remain practical. Additionally, it is unclear how a fundamentally bitmap-based model such as this would scale to much larger numbers of cells. Perhaps vector-drawn “cells”, with internal environments, would be apropos in a future implementation.
Future expansion upon this model could additionally incorporate real-time rendering rather than precomputation of lattices for plotting. In Go, this is conceivable via rendering into two images at alternate timepoints, and plotting one of them, such that one image is always being “painted” and one displayed. This could be further optimized via use of the additional CPU cores (or a GPU) if local rendering is required.
In conclusion, this model, while simple, is an arguably  powerful multifaceted tool for simulating cancerous cell growth in silico.

 
§5. Works Cited
Doyle, Barry, Karol Miller, Adam Wittek, and Poul M.F. Nielsen, eds. Computational Biomechanics for Medicine. New York, NY: Springer New York, 2014. https://doi.org/10.1007/978-1-4939-0745-8.
Ghaemi, Mehrdad, and Amene Shahrokhi. “Combination of the Cellular Potts Model and Lattice Gas Cellular Automata for Simulating the Avascular Cancer Growth.” In Cellular Automata, edited by Samira El Yacoubi, Bastien Chopard, and Stefania Bandini, 297–303. Berlin, Heidelberg: Springer Berlin Heidelberg, 2006.
“Glauber’s Dynamics | Bit-Player.” Accessed November 29, 2019. http://bit-player.org/2019/glaubers-dynamics.
Graner, null, and null Glazier. “Simulation of Biological Cell Sorting Using a Two-Dimensional Extended Potts Model.” Physical Review Letters 69, no. 13 (September 28, 1992): 2013–16. https://doi.org/10.1103/PhysRevLett.69.2013.
Roussos, Evanthia T., John S. Condeelis, and Antonia Patsialou. “Chemotaxis in Cancer.” Nature Reviews. Cancer 11, no. 8 (July 22, 2011): 573–87. https://doi.org/10.1038/nrc3078.
Wolf-Gladrow, Dieter A. Lattice-Gas Cellular Automata and Lattice Boltzmann Models An Introduction. 1st ed. 2000. Lecture Notes in Mathematics, 1725. Berlin, Heidelberg: Springer Berlin Heidelberg, 2000. https://doi.org/10.1007/b72010.
Wolfram, Stephen. “Statistical Mechanics of Cellular Automata.” Reviews of Modern Physics 55, no. 3 (July 1, 1983): 601–44. https://doi.org/10.1103/RevModPhys.55.601.

