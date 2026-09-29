> 🇫🇷 **[Version Française](README.md)** | 🇬🇧 **[English Version](README-EN.md)** | 🎬 **[Step-by-Step Demo Playbook (DEMO_PLAYBOOK-EN.md)](docs/DEMO_PLAYBOOK-EN.md)** | 🚀 **[Live Demo](https://rag.hoffmannw.demo.altostrat.com)** | 🗺️ **[Dendrite v3.3 & Draw.io Diagrams](docs/architecture-gcpdraw.md)**

> 🗺️ **Executive Architecture Diagrams (Dendrite v3.3 Bento 16:9 & Draw.io)**: [**🎨 Launch in Dendrite Studio (1-Click) ↗**](https://dendrite-758054785671.cr.gclb.goog/#code=dGhlbWU6ICJnY3AtYXJjaGl0ZWN0dXJlIgpyZW5kZXJPcmRlcjogbm9kZXMtZmlyc3QKZGlyZWN0aW9uOiByaWdodApzcGFjaW5nOiAyNgoKY29uc3QgR2NwQmx1ZSA9ICIjMWE3M2U4Igpjb25zdCBHY3BHcmVlbiA9ICIjMWU4ZTNlIgpjb25zdCBEYXJrU2xhdGUgPSAiIzIwMjEyNCIKY29uc3QgU3ViVGV4dCA9ICIjNWY2MzY4Igpjb25zdCBDYXJkQm9yZGVyID0gIiNkYWRjZTAiCmNvbnN0IFN1cmZhY2VXaGl0ZSA9ICIjZmZmZmZmIgoKU3R5bGUgQFByb2R1Y3RDYXJkIHsKICB3aWR0aDogMTk2LCBoZWlnaHQ6IDYwLAogIGZpbGw6ICRTdXJmYWNlV2hpdGUsIHN0cm9rZUNvbG9yOiAkQ2FyZEJvcmRlciwgc3Ryb2tlV2lkdGg6IDEsIGJvcmRlclJhZGl1czogOCwKICBmb250Q29sb3I6ICREYXJrU2xhdGUsIHN1YkZvbnRDb2xvcjogJFN1YlRleHQsCiAgZm9udFNpemU6IDEzLCBzdWJGb250U2l6ZTogMTAuNSwgbGFiZWxXZWlnaHQ6IGJvbGQsCiAgdGV4dEFsaWduOiAibGVmdCIsIHRleHRWQWxpZ246ICJtaWRkbGUiLAogIGljb25Qb3NpdGlvbjogImxlZnQiLCBpY29uU2l6ZTogMjYsIHBhZGRpbmc6IDEwLCBzaGFkb3c6IHRydWUKfQpTdHlsZSBAWm9uZUJsdWUgICB7IGZpbGw6ICIjZThmMGZlIiwgc3Ryb2tlOiAiIzViOWJmMyIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZUdyZWVuICB7IGZpbGw6ICIjZTZmNGVhIiwgc3Ryb2tlOiAiIzY4Yjg4ZSIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVBlYWNoICB7IGZpbGw6ICIjZmNlOGU2Iiwgc3Ryb2tlOiAiI2YyOGI4MiIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVB1cnBsZSB7IGZpbGw6ICIjZjNlOGZkIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGJvcmRlclJhZGl1czogMTIgfQoKWm9uZSBAUkFHX0NvbXBhcmlzb25fQXJjaGl0ZWN0dXJlIHsKICB0aXRsZTogIkdDUCBBSSBGb3VuZGF0aW9uIOKAlCBIeWJyaWQgUkFHICYgNC1BZ2VudCBDUkFHIFN3YXJtIENvbXBhcmF0b3IiCiAgc3VidGl0bGU6ICJDbG91ZCBSdW4gdjIgKEdvIDEuMjQgWmVyby1DVkUpIOKAoiBSZWFsLVRpbWUgU1NFIFN0cmVhbWluZyDigKIgVmVydGV4IEFJIEF1dG9yYXRlciAoZXVyb3BlLXdlc3QxKSIKICBpY29uOiAiR29vZ2xlQ2xvdWQiLCBpY29uU2l6ZTogMjQKICBsYXlvdXQ6IG1hdHJpeCwgZ2FwOiAyNCwgcGFkZGluZzogMjQsIGFsaWduOiAic3RyZXRjaCIKICBmaWxsOiAiI2Y4ZmFmZCIsIHN0cm9rZTogIiNkYWRjZTAiLCBjb3JuZXJSYWRpdXM6IDE0CiAgYXJlYXM6IFsKICAgICJ1cHN0cmVhbSAgYXJjaEEgIGFyY2hCIiwKICAgICJ1cHN0cmVhbSAgc3RvcmUgIHN0b3JlIgogIF0KCiAgWm9uZSBAVXBzdHJlYW1fRWRnZSB7CiAgICB0aXRsZTogIjEuIFplcm8tVHJ1c3QgSW5ncmVzcyIKICAgIHN1YnRpdGxlOiAiQ2xvdWQgUnVuIHYyIFNTRSIKICAgIHN0eWxlOiBAWm9uZVBlYWNoLCBhcmVhOiAidXBzdHJlYW0iCiAgICBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkNsb3VkQXJtb3IiLCBpY29uU2l6ZTogMjAKICAgIFthcm1vcl9pYXA6ICJDbG91ZCBBcm1vciAmIElBUCIgfCAiV0FGIEw3ICsgWmVyby1UcnVzdCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkQXJtb3IiIH0KICAgIFtjbG91ZF9ydW5fZ286ICJDbG91ZCBSdW4gdjIgKEdvKSIgfCAiU1NFIFJvdXRlciAmIFgtUmF5IFVJIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwQ2xvdWRSdW4iIH0KICAgIFtkb2NfaW5nZXN0b3I6ICJEb2N1bWVudCBJbmdlc3RvciIgfCAiUERGIENNYXAgLyBET0NYIC8gVFhUIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogIH0KCiAgWm9uZSBAUGF0aF9BX0FnZW50aWNfQ1JBRyB7CiAgICB0aXRsZTogIjJBLiBNb2RlIEFnZW50aWMgUkFHICg0LVN1YmFnZW50IENSQUcgU3dhcm0pIgogICAgc3VidGl0bGU6ICJNdWx0aS1Ib3AgRGVjb21wb3NpdGlvbiAmIFNlbGYtQ29ycmVjdGlvbiIKICAgIHN0eWxlOiBAWm9uZVB1cnBsZSwgYXJlYTogImFyY2hBIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiR29vZ2xlQWdlbnRzIiwgaWNvblNpemU6IDIwCgogICAgWm9uZSBAQ1JBR19QbGFuX1JldHJpZXZlIHsKICAgICAgdGl0bGU6ICJQbGFuICYgTXVsdGktSG9wIFJldHJpZXZlIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJTcGFya2xlcyIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfcGxhbm5lcjogIjEuIFF1ZXJ5UGxhbm5lckFnZW50IiB8ICJEZWNvbXBvc2VzIFN1Yi1RdWVyaWVzIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR29vZ2xlQWdlbnRzIiB9CiAgICAgIFthZ2VudF9yZXRyaWV2ZXI6ICIyLiBIeWJyaWRSZXRyaWV2ZXIiIHwgIjAuNyBDb3NpbmUgKyAwLjMgQk0yNSJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlNlYXJjaCIgfQogICAgfQoKICAgIFpvbmUgQENSQUdfR3JhZGVfU3ludGhlc2l6ZSB7CiAgICAgIHRpdGxlOiAiQ1JBRyBDcml0aWMgJiBTeW50aGVzaXMiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjYTE0MmY0IiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIlZlcnRleEFJIiwgaWNvblNpemU6IDE4CiAgICAgIFthZ2VudF9ncmFkZXI6ICIzLiBHcmFkZXJDcml0aWNBZ2VudCIgfCAiRmlsdGVycyBOb2lzZSBDaHVua3MiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJTaGllbGRDaGVjayIgfQogICAgICBbYWdlbnRfc3ludGg6ICI0LiBDaXRhdGlvblN5bnRoZXNpemVyIiB8ICJHcm91bmRlZCBTU0UgU3RyZWFtIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2VtaW5pIiB9CiAgICB9CiAgfQoKICBab25lIEBQYXRoX0JfQXJlbmFfQW5kX0V2YWwgewogICAgdGl0bGU6ICIyQi4gU3RhbmRhcmQgUkFHLCBNb2RlbCBBcmVuYSAmIEF1dG9yYXRlciBKdWRnZSIKICAgIHN1YnRpdGxlOiAiRmxhc2ggdnMuIFBybyBCZW5jaG1hcmsgJiBMTE0tYXMtYS1KdWRnZSAoLzUpIgogICAgc3R5bGU6IEBab25lQmx1ZSwgYXJlYTogImFyY2hCIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiVmVydGV4QUkiLCBpY29uU2l6ZTogMjAKCiAgICBab25lIEBNdWx0aU1vZGVsX0FyZW5hIHsKICAgICAgdGl0bGU6ICJNdWx0aS1Nb2RlbCBBcmVuYSIsIGxheW91dDogY29sdW1uLCBnYXA6IDE2LCBwYWRkaW5nOiAxNCwgYWxpZ246ICJjZW50ZXIiLAogICAgICBmaWxsOiAiI2ZmZmZmZiIsIHN0cm9rZTogIiM1YjliZjMiLCBjb3JuZXJSYWRpdXM6IDEwLCBpY29uOiAiR2VtaW5pIiwgaWNvblNpemU6IDE4CiAgICAgIFtnZW1pbmlfZmxhc2g6ICJHZW1pbmkgMy41IC8gMy44IEZsYXNoIiB8ICJTdWItNDAwbXMgVFRGVCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdlbWluaSIgfQogICAgICBbZ2VtaW5pX3BybzogIkdlbWluaSAzLjEgUHJvIiB8ICJEZWVwIEZyb250aWVyIFJlYXNvbmluZyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlZlcnRleEFJIiB9CiAgICB9CgogICAgWm9uZSBAQXV0b3JhdGVyX0FuZF9GaW5vcHMgewogICAgICB0aXRsZTogIlF1YWxpdHkgJiBGaW5PcHMgQXVkaXQiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjNWI5YmYzIiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIkFjdGl2aXR5IiwgaWNvblNpemU6IDE4CiAgICAgIFt2ZXJ0ZXhfYXV0b3JhdGVyOiAiVmVydGV4IEFJIEF1dG9yYXRlciIgfCAiR3JvdW5kZWRuZXNzICYgUmVsZXZhbmNlIC81Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiU2VjdXJpdHlDb21tYW5kQ2VudGVyIiB9CiAgICAgIFtmaW5vcHNfYmFkZ2U6ICJSZWFsLVRpbWUgRmluT3BzIiB8ICJUb2tlbiBDb3N0ICh+JDAuMDAwMTEpIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRCaWxsaW5nIiB9CiAgICB9CiAgfQoKICBab25lIEBTaGFyZWRfRGF0YV9Gb3VuZGF0aW9uIHsKICAgIHRpdGxlOiAiMy4gU2hhcmVkIEh5YnJpZCBWZWN0b3IgSW5kZXggJiBDbG91ZCBTdG9yYWdlIFBlcnNpc3RlbmNlIChEaXJlY3QgVlBDIEVncmVzcykiCiAgICBzdWJ0aXRsZTogIkluLU1lbW9yeSA3NjgtZGltIENvc2luZSArIEJNMjUgSW5kZXggQmFja2VkIGJ5IENsb3VkIFN0b3JhZ2UgVUJMQSIKICAgIHN0eWxlOiBAWm9uZUdyZWVuLCBhcmVhOiAic3RvcmUiCiAgICBsYXlvdXQ6IHJvdywgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkdjcERhdGFiYXNlIiwgaWNvblNpemU6IDIwCiAgICBbZ2NzX2NvcnB1c19idWNrZXQ6ICJDbG91ZCBTdG9yYWdlIFVCTEEiIHwgIkNvcnB1cyAmIEpTT05MIEluZGV4Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogICAgW3ZlcnRleF9lbWJlZGRpbmdzOiAiVmVydGV4IEVtYmVkZGluZ3MiIHwgInRleHQtZW1iZWRkaW5nLTAwNCAoNzY4ZCkiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJBSVBsYXRmb3JtIiB9CiAgICBbaHlicmlkX2VuZ2luZTogIkdvIEh5YnJpZCBFbmdpbmUiIHwgIkNvc2luZSAoMC43KSArIEJNMjUgKDAuMykiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJDcHUiIH0KICAgIFtjbG91ZF9sb2dnaW5nX2V2YWw6ICJDbG91ZCBMb2dnaW5nICYgVHJhY2UiIHwgIlNTRSAmIEV2YWwgVGVsZW1ldHJ5Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRMb2dnaW5nIiB9CiAgfQp9CgpbYXJtb3JfaWFwXSAtPiBbY2xvdWRfcnVuX2dvXSB7IHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIxIiwgYmFkZ2VGaWxsOiAkR2NwQmx1ZSwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbY2xvdWRfcnVuX2dvXSAtPiBbZG9jX2luZ2VzdG9yXSB7IHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbY2xvdWRfcnVuX2dvXSAtPiBbYWdlbnRfcGxhbm5lcl0geyBsYWJlbDogIkFnZW50aWMgU1NFIiwgc291cmNlQW5jaG9yOiAicmlnaHQiLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiwgc2VxdWVuY2VCYWRnZTogIjIiLCBiYWRnZUZpbGw6ICRHY3BCbHVlLCBiYWRnZUZvbnRDb2xvcjogJFN1cmZhY2VXaGl0ZSB9ClthZ2VudF9wbGFubmVyXSAtPiBbYWdlbnRfcmV0cmlldmVyXSB7IGxhYmVsOiAiU3ViLVF1ZXJpZXMiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2FnZW50X3JldHJpZXZlcl0gLT4gW2FnZW50X2dyYWRlcl0geyBsYWJlbDogIlRvcC1LIENodW5rcyIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfZ3JhZGVyXSAtPiBbYWdlbnRfc3ludGhdIHsgbGFiZWw6ICJWZXJpZmllZCIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQoKW2FnZW50X3N5bnRoXSAtPiBbZ2VtaW5pX2ZsYXNoXSB7IGxhYmVsOiAiQ29tcGFyZSIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDEuOCwgZGFzaGVkOiB0cnVlLCBjdXJ2ZTogInN0ZXAiIH0KW2dlbWluaV9mbGFzaF0gLT4gW2dlbWluaV9wcm9dIHsgbGFiZWw6ICJBcmVuYSIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDEuOCwgY3VydmU6ICJzdGVwIiB9CltnZW1pbmlfZmxhc2hdIC0+IFt2ZXJ0ZXhfYXV0b3JhdGVyXSB7IGxhYmVsOiAiSnVkZ2UgLzUiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiwgc2VxdWVuY2VCYWRnZTogIjMiLCBiYWRnZUZpbGw6ICRHY3BHcmVlbiwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbdmVydGV4X2F1dG9yYXRlcl0gLT4gW2Zpbm9wc19iYWRnZV0geyBsYWJlbDogIlRva2VuIENvc3QiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMS44LCBjdXJ2ZTogInN0ZXAiIH0KCltkb2NfaW5nZXN0b3JdIC0+IFtnY3NfY29ycHVzX2J1Y2tldF0geyBsYWJlbDogIlBlcnNpc3QgSlNPTkwiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9yZXRyaWV2ZXJdIC0+IFt2ZXJ0ZXhfZW1iZWRkaW5nc10geyBsYWJlbDogIjc2OGQgVmVjdG9yIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9ncmFkZXJdIC0+IFtoeWJyaWRfZW5naW5lXSB7IGxhYmVsOiAiQk0yNSArIENvcyIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbZmlub3BzX2JhZGdlXSAtPiBbY2xvdWRfbG9nZ2luZ19ldmFsXSB7IGxhYmVsOiAiQXVkaXQgTG9ncyIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwR3JlZW4sIHN0cm9rZVdpZHRoOiAxLjgsIGRhc2hlZDogdHJ1ZSwgY3VydmU6ICJzdGVwIiB9Cg==) · [`docs/diagrams/rag_comparison_crag_v3.dendrite`](docs/diagrams/rag_comparison_crag_v3.dendrite) · [`docs/diagrams/rag_comparison_crag_v3.drawio`](docs/diagrams/rag_comparison_crag_v3.drawio) · [`docs/architecture-gcpdraw.md`](docs/architecture-gcpdraw.md)

# RAG Comparison Demo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Vertex_AI-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/vertex-ai)
[![Gemini Models](https://img.shields.io/badge/Gemini-3.5_Flash_%7C_3.8_Flash_%7C_3.1_Pro-8E75B2?style=flat&logo=google-gemini)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Technical demonstration web application comparing keyword-based lexical search with grounded Retrieval-Augmented Generation (RAG) using the Google Gemini 3 model family on Vertex AI.

The application runs as a standalone Go binary compiled with zero external dependencies (`0 dependencies`, `0 CVEs`). It is deployed on **Google Cloud Run** or **Google Kubernetes Engine (GKE) Autopilot** as part of the [GCP AI Foundation Blueprint](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint) architecture.

---

## Features

### 1. Hybrid Search (Semantic & Lexical)
- **Dense embeddings**: Vectorization using the Vertex AI `text-multilingual-embedding-002` model (768-dimensional vectors).
- **Native Go cosine similarity**: High-performance in-memory vector comparison without external vector databases or CGO dependencies.
- **Combined ranking**: Weighted scoring formula $0.65 \times \text{CosineSimilarity} + 0.35 \times \text{LexicalScore}$ (normalized BM25).
- **Automatic fallback**: If the Vertex AI Embedding API is unreachable or rate-limited, the engine automatically falls back to full lexical search with linguistic normalization (accent folding, French elision handling, stop words filtering).

### 2. Gemini 3 Model Suite & Multi-Model Arena
- **Supported models**:
  - `gemini-3.5-flash` (default): Ultra-low latency, high throughput, and native multimodal reasoning.
  - `gemini-3.8-flash`: Frontier agentic capabilities and code-level technical reasoning.
  - `gemini-3.1-pro-preview`: Deep reasoning model for complex architectural analysis and dense synthesis.
  - `gemini-3.5-flash-lite`: High-frequency, cost-sensitive operational workloads.
- **Dynamic model switching**: Hot-swap active models via `/api/model/switch` without container restarts.
- **Display modes (5 interactive views)**:
  - **Standard RAG**: Unified grounded response with retrieved source cards and telemetry.
  - **Split View**: Direct side-by-side comparison between raw lexical excerpts and grounded RAG synthesis.
  - **Model Arena (Dual-Model)**: Concurrent side-by-side inference with two distinct Gemini models on the exact same query.
  - **Triple View**: Three-column comparison (Model A, Model B, Classical Lexical Search).
  - **🤖 Agentic RAG (4-Subagent CRAG Swarm)**: Self-correcting multi-agent pipeline (`QueryPlanner` ➔ `HybridRetriever` ➔ `GraderCritic` with *Auto-Heal* loop ➔ `CitationSynthesizer`) paired with a live **Trace Live** SSE drawer (`event: agent_step`).
- **Real-time telemetry**: In-stream tracking of time-to-first-token (TTFT in ms), total generation latency, and exact token counts parsed from Vertex AI responses.

### 3. On-Demand GenAI Evaluation (Vertex AI Autorater)
- **Unit-level evaluation**: Interactive evaluation button under every generated model response.
- **Calculated metrics**:
  - **Groundedness**: Measures factual alignment between the generated response and the retrieved document chunks (rating from 1 to 5).
  - **QA Relevance**: Evaluates how well the response directly addresses the user's query (rating from 1 to 5).
- **Inspection reports**: Interactive modal displaying numerical ratings, textual rationale produced by the Vertex AI Autorater, and exact contextual excerpts evaluated.
- **Heuristic fallback**: In offline or sandbox environments without Rapid Evaluation quotas, a deterministic n-gram overlap heuristic provides continuous feedback.

### 4. Cloud Storage Persistence & Horizontal Scaling
- **Serialized index**: Uploaded documents and their embeddings are persisted as JSON snapshots in `gs://{GCS_RAG_BUCKET}/index/corpus.json`.
- **Multi-instance support**: Cloud Run instances scale horizontally (`--max-instances=5`) with consistent corpus synchronization.
- **Instant cold starts**: Snapshots are restored from GCS in under one second during instance startup.

---

## Architecture

Review the complete architectural diagram in GCP Draw format in [docs/architecture-gcpdraw.md](docs/architecture-gcpdraw.md).

```
                            [ Client Browser ]
                                     │
                   HTTPS / Identity-Aware Proxy (IAP)
                                     │
                                     ▼
                     [ Google Cloud Load Balancer ]
                                     │
                                     ▼
                     [ Cloud Run : rag-comparison ]
                     ┌────────────────────────────┐
                     │ • Native Go HTTP Server    │
                     │ • Hybrid Search Engine     │
                     │ • 4-Subagent CRAG Swarm    │
                     │ • Static Assets (embed.FS) │
                     └──────┬──────────────┬──────┘
                            │              │
           Vectorization &  │              │  Corpus Storage
           LLM Inference    │              │  & Snapshots
                            ▼              ▼
                    [ Vertex AI ]    [ Cloud Storage ]
                    • Gemini 3.5/3.8 • gs://{BUCKET}/index/
                    • Embeddings 002
                    • Rapid Eval
```

---

## REST API Reference

All endpoints are served by the Go HTTP server on the configured port (`PORT`, default `8080`).

| Method | Endpoint | Parameters | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/documents` | None | Returns the list of all indexed documents and metadata (title, size, chunks). |
| `POST` | `/api/documents/upload` | Multipart form (`files[]`) | Uploads and indexes new documents (PDF, text, markdown) with embeddings. |
| `DELETE` | `/api/documents?id={id}` | `id` (document ID) | Deletes a document from the corpus and updates the Cloud Storage snapshot. |
| `DELETE` | `/api/documents?all=true` | `all=true` | Purges all documents from memory and removes the Cloud Storage snapshot. |
| `GET` | `/api/search/classic` | `q={query}` | Performs keyword-based lexical search and returns raw snippet excerpts. |
| `GET` | `/api/chat/stream` | `q={query}`, `model={model_id}`, `mode=agentic` *(opt.)* | Streams grounded RAG or 4-subagent CRAG responses via Server-Sent Events (`event: agent_step`, `retrieval`, `token`, `metrics`). |
| `GET` | `/api/models` | None | Lists available Gemini models, specifications, and the active default model. |
| `POST` | `/api/model/switch` | JSON body `{"model": "id"}` | Dynamically changes the active default Gemini model. |
| `POST` | `/api/evaluate` | JSON body `{"query", "prediction", "context"}` | Runs an evaluation for Groundedness and QA Relevance via Vertex AI Rapid Evaluation. |
| `GET` | `/healthz` | None | Liveness and readiness health check probe for Cloud Run and Kubernetes. |

---

## Environment Variables

Configure application behavior using the following environment variables:

| Variable | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `GCP_PROJECT` (or `PROJECT_ID`) | String | Auto-detected via GCP metadata | Google Cloud project ID hosting Vertex AI. |
| `GCP_REGION` (or `REGION`) | String | `europe-west1` | Target Google Cloud region for regional Vertex AI calls. |
| `GEMINI_MODEL` | String | `gemini-3.5-flash` | Default Gemini model configured at startup. |
| `GCS_RAG_BUCKET` | String | None (optional) | Cloud Storage bucket name for corpus snapshot persistence (`corpus.json`). |
| `PORT` | Integer | `8080` | TCP port on which the HTTP server listens. |

---

## Local Development

### Prerequisites
- Go 1.25 or higher installed.
- Google Cloud SDK (`gcloud`) authenticated with access to a Google Cloud project.
- Required IAM roles: `roles/aiplatform.user` and `roles/storage.objectAdmin`.

### Run Locally

```bash
# 1. Navigate to the source directory
cd src

# 2. Set environment variables (GCP_PROJECT or PROJECT_ID)
export GCP_PROJECT=$(gcloud config get-value project)
export GCP_REGION="europe-west1"
export GEMINI_MODEL="gemini-3.5-flash"
export GCS_RAG_BUCKET="${GCP_PROJECT}-rag-docs"

# 3. Run unit tests
go test -v ./...

# 4. Start the server
go run main.go
```

Access the web interface at `http://localhost:8080`.

---

## Deployment

Refer to the [Deployment Guide](docs/DEPLOYMENT_GUIDE.md) for step-by-step production deployment procedures.

### Option 1: Google Cloud Run (Serverless)

```bash
./deploy/cloudrun/deploy.sh
```

### Option 2: GKE Autopilot (GCP AI Foundation Blueprint)

```bash
./scripts/deploy-to-blueprint.sh --blueprint-dir=../gcp-ai-foundation-blueprint
```

---

## Security & Governance

- **Zero hardcoded credentials**: The container image contains no service account keys or static secrets. Authentication is managed exclusively through **Workload Identity** (GKE) and Cloud Run execution service accounts via instance metadata (`http://metadata.google.internal`).
- **Network perimeter protection**: Recommended deployment with `--ingress=internal-and-cloud-load-balancing`, protected by Cloud Armor WAF and Identity-Aware Proxy (IAP).
- **Zero third-party runtime dependencies**: The Go backend imports only the Go standard library (`0 direct dependencies`, `0 CVEs` in container vulnerability scans).

---

## 🤖 Dual-Layer Agentic Architecture (Runtime CRAG Swarm & M1L1 Skills)

This repository implements a **two-tier complementary Agentic AI architecture**:
- 🚀 **Layer 2 (Production Run-Time)**: A **4-Subagent Corrective RAG (CRAG) Swarm** embedded natively inside the Go binary (`src/main.go`), triggered in real time by end users from the web UI.
- 🛠️ **Layer 1 (Engineering Build-Time)**: **2 Specialized Subagents and 1 M1L1 Skill** (`.agents/`), triggered inside the IDE/CLI during code development and before every `git commit`.

### 🔄 Sequence Diagram: How the 4 CRAG Subagents Enter into Action Live

When a user clicks the **`🤖 Agentic RAG`** button in the web interface (`https://rag.hoffmannw.demo.altostrat.com`) and submits a question, the backend streams each subagent's execution live over SSE (`event: agent_step`) into the **Trace Live** drawer:

```mermaid
sequenceDiagram
    autonumber
    actor User as 👤 End User (Web UI)
    participant UI as 🖥️ Trace Live Drawer (SSE)
    participant A1 as 🧠 1. QueryPlannerAgent
    participant A2 as 🔍 2. HybridRetrieverAgent
    participant A3 as ⚖️ 3. GraderCriticAgent (CRAG)
    participant A4 as ✍️ 4. CitationSynthesizerAgent

    User->>UI: Toggles "🤖 Agentic RAG" & submits question
    UI->>A1: GET /api/chat/stream?mode=agentic&q=...
    A1-->>UI: SSE agent_step (status: done, 3 sub-queries planned)
    A1->>A2: Dispatches lexical + semantic sub-queries
    A2->>A2: Parallel Dense Cosine (768d) + BM25 search (RRF k=60)
    A2-->>UI: SSE agent_step (status: done, N deduplicated chunks)
    A2->>A3: Submits candidate chunks for factual grading
    alt Insufficient Coverage (Score < 7/10) — Self-Healing Loop
        A3-->>UI: SSE agent_step (status: heal, query rewrite + expanded Top-K)
        A3->>A2: Re-runs HybridRetriever with rewritten query
        A2-->>A3: Enriched document chunks
    else Sufficient Coverage (Score >= 7/10)
        A3-->>UI: SSE agent_step (status: done, relevance verified)
    end
    A3->>A4: Passes verified context chunks
    A4-->>UI: Streams SSE tokens + inline citations [Doc, Chunk #X]
```

### 📊 Summary Matrix: Where and How Each Agent Operates

| Agent / Skill | Layer | Where does it live? | How / When does it enter into action? | Role & Added Value |
| :--- | :--- | :--- | :--- | :--- |
| **`QueryPlannerAgent`** | **Layer 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Step 1** as soon as a query is submitted in `🤖 Agentic RAG` mode. | Decomposes complex, multi-faceted user questions into targeted lexical and semantic sub-queries. |
| **`HybridRetrieverAgent`** | **Layer 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Step 2** after query planning (and re-invoked during *Auto-Heal*). | Executes parallel **Dense Cosine + Sparse BM25 with Reciprocal Rank Fusion ($k=60$)** and deduplicates chunks. |
| **`GraderCriticAgent`** | **Layer 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Step 3** prior to any final LLM generation. | Grades chunk relevance (`/10`). Automatically triggers a **Self-Correction (Query Rewrite)** loop if coverage is insufficient. |
| **`CitationSynthesizerAgent`** | **Layer 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Step 4** once chunks are certified by `GraderCriticAgent`. | Streams the grounded synthesis over SSE with inline citations `[Doc: <Title>, Chunk #X]` and triggers Vertex AI Autorater evaluation. |
| **[`rag-eval-scientist`](.agents/agents/rag-eval-scientist.md)** | **Layer 1** *(Build-Time)* | `.agents/agents/rag-eval-scientist.md` | Inside **Jetski / Antigravity / Gemini CLI** when tuning RRF/BM25 weights or evaluation prompts. | Audits **RRF ($k=60$)**, BM25 ($k_1=1.2, b=0.75$), and **LLM-as-a-Judge** faithfulness/relevance metrics. |
| **[`go-concurrency-reviewer`](.agents/agents/go-concurrency-reviewer.md)** | **Layer 1** *(Build-Time)* | `.agents/agents/go-concurrency-reviewer.md` | Inside **Jetski / Antigravity / Gemini CLI** before committing Go backend changes (`src/main.go`). | Audits `sync.RWMutex` lock discipline, SSE `r.Context().Done()` goroutine cleanup, and non-blocking GCS index persistence. |
| **[`rag-benchmark-and-ci`](.agents/skills/rag-benchmark-and-ci/SKILL.md)** | **Layer 1** *(Gatekeeper)* | `.agents/skills/rag-benchmark-and-ci/scripts/verify.sh` | Executed in the terminal before every `git commit` or Cloud Run deployment. | Runs `go vet ./...`, `go test -v -race ./...`, enforces **Zero External Dependencies** (`src/go.mod`), and verifies the 4 CRAG agents. |

### 🎬 3-Minute Customer Demo Playbook (CE Walkthrough)

1. **Step 1 — Trigger the CRAG Swarm Live in the Browser (Run-Time)**:
   - Open **[RAG Comparison Demo](https://rag.hoffmannw.demo.altostrat.com)** and click the **`🤖 Agentic RAG`** button in the top navigation bar.
   - Copy-paste a multi-faceted technical prompt:
     > `"Compare the Zero-Trust security perimeter (IAP, Cloud Armor) with the WORM backup strategy and explain how hybrid RAG prevents hallucinations."`
   - **What to highlight on screen**:
     - The **Trace Live** drawer expands automatically above the answer, streaming the 4 subagent cards in real time (`QueryPlanner` ➔ `HybridRetriever` ➔ `GraderCritic` ➔ `CitationSynthesizer`) with per-agent millisecond latency and generated sub-queries.
     - Click **`⭐ Evaluate (Vertex AI)`** beneath the response to display the **Groundedness (`/5`)** and **QA Relevance (`/5`)** scores.

2. **Step 2 — Inspect the Raw `event: agent_step` SSE Stream via `curl`**:
   ```bash
   curl -N "https://rag.hoffmannw.demo.altostrat.com/api/chat/stream?mode=agentic&q=Cloud+Armor+and+Hybrid+RRF"
   ```
   *(Streams live `event: agent_step` JSON payloads followed by synthesized response tokens).*

3. **Step 3 — Showcase the Build-Time Engineering Subagents & M1L1 Gatekeeper**:
   - In **Jetski / Antigravity / Gemini CLI**, copy-paste:
     > `"Invoke rag-eval-scientist to audit the Reciprocal Rank Fusion (k=60) scoring and GraderCriticAgent self-correction threshold in src/main.go."`
   - Run the M1L1 gatekeeper script:
     ```bash
     ./.agents/skills/rag-benchmark-and-ci/scripts/verify.sh
     ```

---

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.



