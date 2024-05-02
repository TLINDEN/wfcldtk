out=images/output.png
in=images/inputtilemap.png
w=8
h=8
cell=100



build:
	go build .

test: build render show

render:
	./wfcldtk $(in) $(w) $(h) $(cell) $(out)

show:
	display $(out)
