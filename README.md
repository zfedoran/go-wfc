# go-wfc
Procedurally-generated tile maps using wave function collapse. 

## Demos

Live demo (wasm):
* https://zfedoran.github.io/go-wfc-example/

Live algorithm animation (wasm):
* https://zfedoran.github.io/go-wfc-algorithm/


## Overview
This package uses the *Wave Function Collapse* algorithm as described by [Oskar
Stålberg](https://www.youtube.com/watch?v=0bcZb-SsnrA&t=350s).

The wave function collapse algorithm is a recursive algorithm that picks a
random tile for a slot on the output image and removes impossible neighbors
until only a single possibility remains. The algorithm is covered in more 
detail [below](#algorithm).

If you need a generic algorithm, then please check out the golang fork of 
the [original work](https://github.com/shawnridgeway/wfc).

<img src="/doc/images/banner.jpg?raw=true" width="80%">

## Why?

There is already a wfc golang library so why another one? The existing one is a
lot more generic and quite a bit slower as a result. Also, the tile setup for
the original implementation can be very tedious. Additionally, I found it hard
to follow and modify. 

This variation follows Oskars work and aims to simplify the original work for
easy integration into games.

## Tiles

You'll need a set of tiles, also referred to as `modules`. These will 
need to be designed to fit together. The tiles should have matching 
colors along edges that should appear next to each other in the final 
output. Other than that, no manual setup or description files are needed.

You should reference the [example folder](/example/) when making a tile set. 

![sample](/doc/images/tiles.png?raw=true)

* Any number of tiles can be used as input, but the recommendation is to start
with a small set until you're comfortable with the constraint system.

* You can create alternate tiles, meaning tiles with the same sets of colors
along all four edges. The probability for each alternate tile is the same. If
you'd like to increase/lower the probability of a particular tile, simply
duplicate the image reference when calling the initialize method.

* Unlike the original WFC implementation, no manual setup or description files
are needed.

## Adjacencies / Constraints

The wave function collapse algorithm requires some kind of adjacency mapping in
order to remove impossible tiles and stitch together a possible output using the
input set.

By default, the package uses the color values along the four edges of each tile
(`Up`, `Down`, `Left`, `Right`) to build constraints. But, you can
[customize](#custom-constraints) this behaviour if you'd like.

<img src="/doc/images/constraints.jpg?raw=true" width="50%">

When designing your tiles, think about how the color values line up. They should
be exactly the same on the middle 3 points for two potentially adjacent tiles.
For example, the following tiles could appear as shown below because they share
the same colors on the bottom of the first and the top of the second.

<img src="/doc/images/adjacencies.jpg?raw=true" width="50%">

You can view the default adjacency constraint implementation
[here](https://github.com/zfedoran/go-wfc/blob/main/pkg/wfc/constraint.go#L36). 
It scans colors along each edge of each input tile. These colors are turned into
a hash that represents that edge. Any tiles that have the same hash value in the
opposite direction are considered possible adjacencies automatically.


## Contradictions

It is possible that the wave can collapse into a state that has a contradiction.
For example, sky beneath the ground. If a contradiction is found, the algorithm
re-tries until the maximum number of attempts is reached.

When exporting an image, if you see a red tile, you've got a contradiction. If
you keep seeing these, your tileset likely has an issue.

<img src="/doc/images/contradiction.png?raw=true" width="50%">

## Results

Here are some example outputs for a 8 x 8 grid.

![results](/doc/images/permutations.jpg?raw=true)


## Quick Start

You'll need to load a set of tiles (images) into an array. A convenience
function is provided by this package but you can use any method you'd like.

```go
  // Load the input tile images (any order and count is fine)
  var input_images []image.Image
  input_images, err = wfc.LoadImageFolder(tileset_folder)
  if err != nil {
    panic(err)
  }
```

Next, initialize a wave function with the desired output size (in units of
tiles). For example, lets say that you want your output image to be 32 x 8
tiles, you'd pass in the following.

```go
  // Setup the initialized state. The output image will be 32 x 8 tiles.
  wave := wfc.New(input_images, 32, 8)
  wave.Initialize(42) // seed: 42
```

Finally, collapse the wave into a single state (if possible).

```go
  // Collapse the wave function (make up to 100 attempts)
  err = wave.Collapse(200)
  if err != nil {
    panic(err)
  }
```

Optionally, you can export the collapsed wave to an image.

```go
  // Lets generate an image
  output_image := wave.ExportImage()
  wfc.SaveImage("wave.png", output_image)
```

Or, you can review the results manually to do custom rendering in your game.

```go
  for _, slot := range wave.PossibilitySpace {
    if len(slot.Superposition) == 1 {
      // successfully collapsed slot
      ...
    }
    if len(slot.Superposition) == 0 {
      // contradiction
      ...
    }
  }
```

## Full example:

```go
import "github.com/zfedoran/go-wfc/pkg/wfc"

func collapseWave(tileset_folder, output_image string) {
  // This is just a `[]image.Image`, you can use whatever loader function you'd like
  images, err := wfc.LoadImageFolder(tileset_folder)
  if err != nil {
    panic(err)
  }

  // The random seed to use when collapsing the wave
  // (given the same seed number, the Collapse() fn would generate the same state every time)
  seed := int(time.Now().UnixNano())

  // Setup the initialized state
  wave := wfc.New(images, 32, 8)
  wave.Initialize(seed)
  
  // Collapse the wave function (make up to 100 attempts)
  err = wave.Collapse(200)
  if err != nil {
    // don't panic here, we want to generate the image anyway
    fmt.Printf("unable to generate: %v", err)
  }

  // Lets generate an image
  output := wave.ExportImage()
  wfc.SaveImage(output_image, output)

  fmt.Printf("Image saved to: %s\n", output_image)
}
```

Complete source can be found here:
[example/main.go](example/main.go)

Also, check out the animated version:
https://github.com/zfedoran/go-wfc-example


## Custom Constraints
If you'd like to customize or change this logic, you are able to pass in a
custom constraint function.

You can choose a different number of lookup points (3 is the default). For
example, 2 lookup points.

```go
  wave.NewWithCustomConstraints(tiles, width, height, wfc.GetConstraintFunc(2))
```

Or, you can provide your own.

```go
  wave.NewWithCustomConstraints(tiles, width, height, 
    func(img image.Image, d Direction) ConstraintId {
    ...
  })
```

## Algorithm

The algorithm is covered in detail here:
https://www.youtube.com/watch?v=0bcZb-SsnrA&t=350s

1) A set of input image tiles (or modules) are loaded into memory. 
2) A wave function is defined with the desired output width and height (in
units of tiles).
3) The wave function is initialized such that each output tile (or slot) is in a
superposition of all provided input tiles.
4) A random slot is selected and collapsed into a random input tile.
5) Each of the neighboring slots is now evaluated to verify if there are any
input tiles that can fit next to the collapsed tile. Any impossible tiles
(or modules) are removed.
6) If the state of any of the neighboring tiles was changed in step 5), then
recurse into it's neighbors to remove impossible tiles.
7) If there are no possible tiles left at any point, a contradiction has been
found and we need to go back to step 3) and try again.
8) Once no more changes are left to propagate, go to step 4) and recurse until
all slots are collapsed to a single state.

Or, you if you prefer, here is the [actual implementation](https://github.com/zfedoran/go-wfc/blob/main/pkg/wfc/wave.go#L109).

## Artwork

The awesome artwork in this repository was done by
[@makionfire](https://twitter.com/makionfire). If you need help designing a tile
set, I highly recommend reaching out to her. A huge shout-out to `@makionfire`
for letting me use this tileset.

The artwork itself does **not** fall under the MIT licence.

## Licence

The licence for the source code in this package is MIT. Meaning, do whatever
you'd like but we'd love a shoutout. The goal is to get more folks to build
games with golang.

If you like this work and want to buy me or the artist a coffee or beer, you're
free to do so by sending to some SOL to
[🍺💵.sol](https://naming.bonfida.org/#/domain/%F0%9F%8D%BA%F0%9F%92%B5)
