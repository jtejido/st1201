# st1201
## Floating Point to Integer Mapping

Many data values which are measured or computed in floating point are inserted into KLV for
transmission between systems. Oftentimes, values do not fully utilize the floating-point range or
precision afforded, and thus there is opportunity to reduce the number of bytes used to represent
the data. Additionally, there are special cases in measurements and computations that need to be
communicated; for example, a sensor value that means “beyond measurement range” or in
computation when a divide by zero occurs (i.e. +infinity). This ST applies to IEEE 754 
floating-point values (all precisions; i.e. 16, 32, 64 and 128 bit), including many of the IEEE
special values of infinity, and NaN’s.

### MISB ST 1201.3

This Standard (ST) describes the method for mapping floating-point values to integer values and
the reverse, mapping integer values back to their original floating-point value to within an
acceptable precision. There are many ways of optimizing the transmission of floating-point
values from one system to another; the purpose of this ST is to provide a single method which
can be used for all floating-point ranges and valid precisions. This ST provides a method for a
forward and reverse linear mapping of a specified range of floating-point values to a specified
integer range of values based on the number of bytes desired to be used for the integer value.
Additionally, it provides a set of special values which can be used to transmit non-numerical
“signals” to a receiving system. 
