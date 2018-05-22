/*
Package spec build the framework needed to abstract features on tests.

This is made initializing specification for test, using spec.New(), and
then calling the different prints, in the right order. It will be
necessary to update information about When, It and Assert fields to
print correct information.

The colors used on the output, can be changed using spec.SetConfig(),
where you can create different Configuration object, using different
colors for each type of line printed.
*/
package spec
