#!/bin/sh

# TODO Move this into go app

set -u -e

egrep -v "^#|^[[:space:]]+$" /mnt/z/Musik/barcodes.csv | while read LINE; do
    #echo "BUILD SVG FOR: ${LINE}" 

    BARCODE="$(echo ${LINE} | awk -F',' '{print $1}')"
    ACTION="$(echo ${LINE} | awk -F',' '{print $2}')"

    echo "${BARCODE}" | barcode -umm -g "50x20" -e code128 -S > ./output/${BARCODE}.svg
    
    # Modify SVG - Remove orginal text
    if [ "${ACTION}" = "mpc_add_and_play" ]; then
      # Get Musicdir and escape chars
      MUSICDIR="$(echo ${LINE} | awk -F',' '{print $3}')"
      MUSICDIR="$(echo "${MUSICDIR}" | sed 's/&/\&amp;/g; s/</\&lt;/g; s/>/\&gt;/g; s/"/\&quot;/g; s/'"'"'/\&#39;/g')"

      sed -i 's#<text.*##' ./output/${BARCODE}.svg
      sed -i 's#</svg>##' ./output/${BARCODE}.svg

      echo "<text x=\"10pt\" y=\"66pt\" fill=\"black\" style=\"font-family:Helvetica;font-size:10pt\">$(dirname "${MUSICDIR}")</text>" >> ./output/${BARCODE}.svg
      echo "<text x=\"10pt\" y=\"76pt\" fill=\"black\" style=\"font-family:Helvetica;font-size:10pt\">$(basename "${MUSICDIR}")</text>" >> ./output/${BARCODE}.svg
      echo "</svg>" >> ./output/${BARCODE}.svg
    fi
done

cat ./01_layout_2_per_page_partA.html

PAGE_ODD_OR_EVEN="odd"

find ./output/ -name "[0-9]*.svg" | sort | xargs -L24 | while read SVG1 SVG2 SVG3 SVG4 SVG5 SVG6 SVG7 SVG8 SVG9 SVG10 SVG11 SVG12 SVG13 SVG14 SVG15 SVG16 SVG17 SVG18 SVG19 SVG20 SVG21 SVG22 SVG23 SVG24; do
  if [ "${PAGE_ODD_OR_EVEN}" = "odd" ]; then
    PAGE_ODD_OR_EVEN="even"
  else
    PAGE_ODD_OR_EVEN="odd"
  fi

cat <<EOF
<div class="mypage ${PAGE_ODD_OR_EVEN}">
  <div class="top">
    <table>
      <tr>
        <td><img class="barcode" src="${SVG1}" /></td>
        <td><img class="barcode" src="${SVG2}" /></td>
        <td><img class="barcode" src="${SVG3}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG4}" /></td>
        <td><img class="barcode" src="${SVG5}" /></td>
        <td><img class="barcode" src="${SVG6}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG7}" /></td>
        <td><img class="barcode" src="${SVG8}" /></td>
        <td><img class="barcode" src="${SVG9}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG10}" /></td>
        <td><img class="barcode" src="${SVG11}" /></td>
        <td><img class="barcode" src="${SVG12}" /></td>
      </tr>
    </table>
  </div>
  <div class="bottom">
    <table>
      <tr>
        <td><img class="barcode" src="${SVG13}" /></td>
        <td><img class="barcode" src="${SVG14}" /></td>
        <td><img class="barcode" src="${SVG15}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG16}" /></td>
        <td><img class="barcode" src="${SVG17}" /></td>
        <td><img class="barcode" src="${SVG18}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG19}" /></td>
        <td><img class="barcode" src="${SVG20}" /></td>
        <td><img class="barcode" src="${SVG21}" /></td>
      </tr>
      <tr>
        <td><img class="barcode" src="${SVG22}" /></td>
        <td><img class="barcode" src="${SVG23}" /></td>
        <td><img class="barcode" src="${SVG24}" /></td>
      </tr>
    </table>
  </div>
</div>
EOF
done

cat ./01_layout_2_per_page_partC.html