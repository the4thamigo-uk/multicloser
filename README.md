# multicloser

_Small go package to manage closing of a set of resources_

The aim of this library is to provide a way to defer any number of Close() functions
to a scope which is not the end of the current function.
