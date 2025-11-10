#!/usr/bin/env python3
"""
Generate minimal valid DCP (DNG Camera Profile) files for testing.

DCP files are TIFF containers with Adobe Camera Profile XML in tag 50740.
This script creates synthet DCP files for Recipe parser testing.
"""

import struct
import xml.etree.ElementTree as ET

def create_minimal_dcp():
    """Create a minimal linear DCP with identity matrices."""

    # Create Camera Profile XML
    profile_xml = '''<?xml version="1.0" encoding="UTF-8"?>
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
                   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <crs:ProfileName>Minimal Linear Test Profile</crs:ProfileName>
  <crs:ToneCurve>
    <rdf:Seq>
      <rdf:li>0, 0</rdf:li>
      <rdf:li>64, 64</rdf:li>
      <rdf:li>128, 128</rdf:li>
      <rdf:li>192, 192</rdf:li>
      <rdf:li>255, 255</rdf:li>
    </rdf:Seq>
  </crs:ToneCurve>
  <crs:ColorMatrix1>
    <rdf:Seq>
      <rdf:li>1.0 0.0 0.0</rdf:li>
      <rdf:li>0.0 1.0 0.0</rdf:li>
      <rdf:li>0.0 0.0 1.0</rdf:li>
    </rdf:Seq>
  </crs:ColorMatrix1>
</crs:CameraProfile>'''

    xml_bytes = profile_xml.encode('utf-8')

    # Build TIFF structure (little-endian)
    tiff = bytearray()

    # TIFF Header
    tiff += b'II'  # Little-endian marker
    tiff += struct.pack('<H', 42)  # TIFF version 42
    tiff += struct.pack('<I', 8)  # Offset to first IFD (after header)

    # Image File Directory (IFD)
    ifd_start = len(tiff)

    # IFD entry count (we'll write 1 entry: tag 50740)
    tiff += struct.pack('<H', 1)  # Number of entries

    # IFD Entry for tag 50740 (CameraProfile)
    # Each entry: tag(2) + type(2) + count(4) + value/offset(4) = 12 bytes
    tiff += struct.pack('<H', 50740)  # Tag ID: CameraProfile
    tiff += struct.pack('<H', 1)      # Type: BYTE
    tiff += struct.pack('<I', len(xml_bytes))  # Count: XML length

    # Value offset (points to XML data after IFD)
    xml_offset = ifd_start + 2 + 12 + 4  # After entry count, entry, and next IFD offset
    tiff += struct.pack('<I', xml_offset)

    # Next IFD offset (0 = no more IFDs)
    tiff += struct.pack('<I', 0)

    # XML data
    tiff += xml_bytes

    return bytes(tiff)

def create_portrait_adjusted_dcp():
    """Create a portrait-style DCP with tone adjustments."""

    # Create Camera Profile XML with S-curve
    profile_xml = '''<?xml version="1.0" encoding="UTF-8"?>
<crs:CameraProfile xmlns:crs="http://ns.adobe.com/camera-raw-settings/1.0/"
                   xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <crs:ProfileName>Portrait Adjusted Test Profile</crs:ProfileName>
  <crs:ToneCurve>
    <rdf:Seq>
      <rdf:li>0, 10</rdf:li>
      <rdf:li>64, 70</rdf:li>
      <rdf:li>128, 148</rdf:li>
      <rdf:li>192, 200</rdf:li>
      <rdf:li>255, 250</rdf:li>
    </rdf:Seq>
  </crs:ToneCurve>
  <crs:ColorMatrix1>
    <rdf:Seq>
      <rdf:li>1.0 0.0 0.0</rdf:li>
      <rdf:li>0.0 1.0 0.0</rdf:li>
      <rdf:li>0.0 0.0 1.0</rdf:li>
    </rdf:Seq>
  </crs:ColorMatrix1>
</crs:CameraProfile>'''

    xml_bytes = profile_xml.encode('utf-8')

    # Build TIFF structure (little-endian)
    tiff = bytearray()

    # TIFF Header
    tiff += b'II'  # Little-endian marker
    tiff += struct.pack('<H', 42)  # TIFF version 42
    tiff += struct.pack('<I', 8)  # Offset to first IFD

    # Image File Directory (IFD)
    ifd_start = len(tiff)

    # IFD entry count
    tiff += struct.pack('<H', 1)  # Number of entries

    # IFD Entry for tag 50740
    tiff += struct.pack('<H', 50740)  # Tag ID
    tiff += struct.pack('<H', 1)      # Type: BYTE
    tiff += struct.pack('<I', len(xml_bytes))  # Count

    # XML offset
    xml_offset = ifd_start + 2 + 12 + 4
    tiff += struct.pack('<I', xml_offset)

    # Next IFD offset
    tiff += struct.pack('<I', 0)

    # XML data
    tiff += xml_bytes

    return bytes(tiff)

if __name__ == '__main__':
    # Generate minimal linear DCP
    print("Generating minimal-linear.dcp...")
    with open('minimal-linear.dcp', 'wb') as f:
        f.write(create_minimal_dcp())
    print("✓ Created minimal-linear.dcp")

    # Generate portrait adjusted DCP
    print("Generating portrait-adjusted.dcp...")
    with open('portrait-adjusted.dcp', 'wb') as f:
        f.write(create_portrait_adjusted_dcp())
    print("✓ Created portrait-adjusted.dcp")

    print("\nDCP test files generated successfully!")
