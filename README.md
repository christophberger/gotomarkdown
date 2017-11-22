> Note: This repo is not maintained anymore. It is superseded by the `gotohugo` repo.

# GoToMarkdown: Convert commented Go files to Markdown (plus an extra gimmick)

## Description

`gotomarkdown` converts a .go file into a Markdown file. Comments can (and should) contain [Markdown](daringfireball.net/projects/markdown) text. Comment delimiters are stripped, and Go code is put into code fences.

Extra: A non-standard "Hype" tag can refer to Tumult Hype HTML animations. This tag is replaced by the corresponding HTML snippet that loads the animation. Create the anmiation from Tumult Hype by exporting to HTML5, with the "Also save HTML file" checkbox checked. `gotomarkdown` can then extract the HTML snippet from the HTML file and can copy the `hyperesources` directory to the output folder.

## Usage

	gotomarkdown [-outdir "path/to/outputDir"] [-nocopy] <gofile.go>

### Flags

*`-outdir`: Specifies the output directory. Defaults to "out".
*`-nocopy`: If set, the image and animation files do not get copied to the output directory. 
*`-subdir`: If set, the image and animation files are copied into &lt;outdir>/&lt;subdir>, rather than into &lt;outdir>.

## License

(c) 2016 Christoph Berger. All Rights Reserved. 
This code is governed by a BSD 3-clause license that can be found in LICENSE.txt.


