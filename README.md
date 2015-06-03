spectrum
========

A library for handling spectrum data

TODO
----

* Tests for arithmetic operations especially when using cubic splines
* Function for efficient X sorting status check. A spectrum must be always sorted by X and the disordered status must be easily understood if X1 > XN. So the sorting should take place only on a spectrum initialization and in case of possible arbitrary X modifications such as ModifyX method.
* Smoothing by Savitsky-Golay and/or Holoborodko

TOTO
----

* Units of X might be incorporated into the Spectrum or SpectrumWrapper (spectool) type.
* Output formatting 
	** Origin
	** Matlab
	** Julia
	** JSON
	** etc.