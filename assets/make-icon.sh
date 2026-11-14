#!/bin/zsh
# Regenerates assets/AppIcon.icns from assets/summon.svg (needs rsvg-convert).
set -e

# Work from the assets folder, with a scratch dir that's cleaned up on exit
cd "${0:A:h}"
tmp=$(mktemp -d)
trap "rm -rf $tmp" EXIT

# Wrap the svg: shrink to 824px, centered on 1024, clipped to macOS-style rounded corners
{
  echo '<svg xmlns="http://www.w3.org/2000/svg" width="1024" height="1024">'
  echo '<clipPath id="r"><rect x="100" y="100" width="824" height="824" rx="185"/></clipPath>'
  echo '<g clip-path="url(#r)">'
  sed 's/<svg [^>]*>/<svg x="100" y="100" width="824" height="824" viewBox="0 0 400 400" fill="none">/' summon.svg
  echo '</g></svg>'
} > $tmp/icon.svg

# Render to a 1024px png
rsvg-convert $tmp/icon.svg -o $tmp/icon.png

# Make every size macOS wants (normal + @2x retina)
mkdir $tmp/AppIcon.iconset
for s in 16 32 128 256 512; do
  sips -z $s $s $tmp/icon.png --out $tmp/AppIcon.iconset/icon_${s}x${s}.png >/dev/null
  sips -z $((s*2)) $((s*2)) $tmp/icon.png --out $tmp/AppIcon.iconset/icon_${s}x${s}@2x.png >/dev/null
done

# Pack them into a single .icns
iconutil -c icns $tmp/AppIcon.iconset -o AppIcon.icns
