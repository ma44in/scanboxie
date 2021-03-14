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

cat <<EOF
<html>
<head>
<style>
  .header, .header-space,
  .footer, .footer-space {
    height: 50px;
  }
  
  .header {
    background: red;
    position: fixed;
    top: 0;
  }
  
  .footer {
    position: fixed;
    bottom: 0;
  }

  @page:first {
    .header {
      content: normal;
    }
  }

  @media print {
    .header {
        background-color: red !important;
        -webkit-print-color-adjust: exact; 
    }
  }

  .center {
    margin: auto;
    width: 100%;
    padding: 10px;
  }

  table.center {
    margin-left:auto; 
    margin-right:auto;
  }

  body {
    font-family: Algerian;
  }

  .barcode-small {
    height: 80px;
  }

</style>
</head>
<body>
EOF

cat ./book_page0.html
cat ./book_page1.html

BARCODE_IMG_ATTRIBUTES=""
#BARCODE_IMG_ATTRIBUTES='height="100px"'

find ./output/ -name "[0-9]*.svg" | sort | xargs -L12 | while read SVG1 SVG2 SVG3 SVG4 SVG5 SVG6 SVG7 SVG8 SVG9 SVG10 SVG11 SVG12; do
cat <<EOF
<div style="page-break-after: always;" class="center">
  <table>
    <thead><tr><td>
      <div class="header-space">&nbsp;</div>
    </td></tr></thead>
    <tbody><tr><td>
      <div class="content">
      
      <table>
        <tr>
          <td width="50%"><img class="barcode" src="${SVG1}"/></td>
          <td width="50%"><img class="barcode" src="${SVG2}"/></td>
        </tr>
        <tr>
          <td><img class="barcode" src="${SVG3}"/></td>
          <td><img class="barcode" src="${SVG4}"/></td>
        </tr>
        <tr>
          <td><img class="barcode" src="${SVG5}"/></td>
          <td><img class="barcode" src="${SVG6}"/></td>
        </tr>
        <tr>
          <td><img class="barcode" src="${SVG7}"/></td>
          <td><img class="barcode" src="${SVG8}"/></td>
        </tr>
        <tr>
          <td><img class="barcode" src="${SVG9}"/></td>
          <td><img class="barcode" src="${SVG10}"/></td>
        </tr>
        <tr>
          <td><img class="barcode" src="${SVG11}"/></td>
          <td><img class="barcode" src="${SVG12}"/></td>
        </tr>
      </table>     
      
      </div>
    </td></tr></tbody>
    <tfoot><tr><td>
      <div class="footer-space">&nbsp;</div>
    </td></tr></tfoot>
  </table>

  <div class='header'>Header</div>
  <div class='footer'>Footer</div>

</div>
EOF
done

cat <<EOF
</body>
</html>
EOF
