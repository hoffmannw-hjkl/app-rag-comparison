> 🇫🇷 **[Version Française](README.md)** | 🇬🇧 **[English Version](README-EN.md)** | 🎬 **[Guide de Démo Pas-à-Pas (DEMO_PLAYBOOK.md)](docs/DEMO_PLAYBOOK.md)** | 🚀 **[Démo Live](https://rag.hoffmannw.demo.altostrat.com)** | 🗺️ **[Schémas Dendrite v3.3 & Draw.io](docs/architecture-gcpdraw.md)**

> 🗺️ **Schémas d'Architecture Exécutifs (Dendrite v3.3 Bento 16:9 & Draw.io)** : [**🎨 Ouvrir dans Dendrite Studio (1-Click) ↗**](https://dendrite-758054785671.cr.gclb.goog/#code=dGhlbWU6ICJnY3AtYXJjaGl0ZWN0dXJlIgpyZW5kZXJPcmRlcjogbm9kZXMtZmlyc3QKZGlyZWN0aW9uOiByaWdodApzcGFjaW5nOiAyNgoKY29uc3QgR2NwQmx1ZSA9ICIjMWE3M2U4Igpjb25zdCBHY3BHcmVlbiA9ICIjMWU4ZTNlIgpjb25zdCBEYXJrU2xhdGUgPSAiIzIwMjEyNCIKY29uc3QgU3ViVGV4dCA9ICIjNWY2MzY4Igpjb25zdCBDYXJkQm9yZGVyID0gIiNkYWRjZTAiCmNvbnN0IFN1cmZhY2VXaGl0ZSA9ICIjZmZmZmZmIgoKU3R5bGUgQFByb2R1Y3RDYXJkIHsKICB3aWR0aDogMTk2LCBoZWlnaHQ6IDYwLAogIGZpbGw6ICRTdXJmYWNlV2hpdGUsIHN0cm9rZUNvbG9yOiAkQ2FyZEJvcmRlciwgc3Ryb2tlV2lkdGg6IDEsIGJvcmRlclJhZGl1czogOCwKICBmb250Q29sb3I6ICREYXJrU2xhdGUsIHN1YkZvbnRDb2xvcjogJFN1YlRleHQsCiAgZm9udFNpemU6IDEzLCBzdWJGb250U2l6ZTogMTAuNSwgbGFiZWxXZWlnaHQ6IGJvbGQsCiAgdGV4dEFsaWduOiAibGVmdCIsIHRleHRWQWxpZ246ICJtaWRkbGUiLAogIGljb25Qb3NpdGlvbjogImxlZnQiLCBpY29uU2l6ZTogMjYsIHBhZGRpbmc6IDEwLCBzaGFkb3c6IHRydWUKfQpTdHlsZSBAWm9uZUJsdWUgICB7IGZpbGw6ICIjZThmMGZlIiwgc3Ryb2tlOiAiIzViOWJmMyIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZUdyZWVuICB7IGZpbGw6ICIjZTZmNGVhIiwgc3Ryb2tlOiAiIzY4Yjg4ZSIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVBlYWNoICB7IGZpbGw6ICIjZmNlOGU2Iiwgc3Ryb2tlOiAiI2YyOGI4MiIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVB1cnBsZSB7IGZpbGw6ICIjZjNlOGZkIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGJvcmRlclJhZGl1czogMTIgfQoKWm9uZSBAUkFHX0NvbXBhcmlzb25fQXJjaGl0ZWN0dXJlIHsKICB0aXRsZTogIkdDUCBBSSBGb3VuZGF0aW9uIOKAlCBIeWJyaWQgUkFHICYgNC1BZ2VudCBDUkFHIFN3YXJtIENvbXBhcmF0b3IiCiAgc3VidGl0bGU6ICJDbG91ZCBSdW4gdjIgKEdvIDEuMjQgWmVyby1DVkUpIOKAoiBSZWFsLVRpbWUgU1NFIFN0cmVhbWluZyDigKIgVmVydGV4IEFJIEF1dG9yYXRlciAoZXVyb3BlLXdlc3QxKSIKICBpY29uOiAiR29vZ2xlQ2xvdWQiLCBpY29uU2l6ZTogMjQKICBsYXlvdXQ6IG1hdHJpeCwgZ2FwOiAyNCwgcGFkZGluZzogMjQsIGFsaWduOiAic3RyZXRjaCIKICBmaWxsOiAiI2Y4ZmFmZCIsIHN0cm9rZTogIiNkYWRjZTAiLCBjb3JuZXJSYWRpdXM6IDE0CiAgYXJlYXM6IFsKICAgICJ1cHN0cmVhbSAgYXJjaEEgIGFyY2hCIiwKICAgICJ1cHN0cmVhbSAgc3RvcmUgIHN0b3JlIgogIF0KCiAgWm9uZSBAVXBzdHJlYW1fRWRnZSB7CiAgICB0aXRsZTogIjEuIFplcm8tVHJ1c3QgSW5ncmVzcyIKICAgIHN1YnRpdGxlOiAiQ2xvdWQgUnVuIHYyIFNTRSIKICAgIHN0eWxlOiBAWm9uZVBlYWNoLCBhcmVhOiAidXBzdHJlYW0iCiAgICBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkNsb3VkQXJtb3IiLCBpY29uU2l6ZTogMjAKICAgIFthcm1vcl9pYXA6ICJDbG91ZCBBcm1vciAmIElBUCIgfCAiV0FGIEw3ICsgWmVyby1UcnVzdCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkQXJtb3IiIH0KICAgIFtjbG91ZF9ydW5fZ286ICJDbG91ZCBSdW4gdjIgKEdvKSIgfCAiU1NFIFJvdXRlciAmIFgtUmF5IFVJIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwQ2xvdWRSdW4iIH0KICAgIFtkb2NfaW5nZXN0b3I6ICJEb2N1bWVudCBJbmdlc3RvciIgfCAiUERGIENNYXAgLyBET0NYIC8gVFhUIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogIH0KCiAgWm9uZSBAUGF0aF9BX0FnZW50aWNfQ1JBRyB7CiAgICB0aXRsZTogIjJBLiBNb2RlIEFnZW50aWMgUkFHICg0LVN1YmFnZW50IENSQUcgU3dhcm0pIgogICAgc3VidGl0bGU6ICJNdWx0aS1Ib3AgRGVjb21wb3NpdGlvbiAmIFNlbGYtQ29ycmVjdGlvbiIKICAgIHN0eWxlOiBAWm9uZVB1cnBsZSwgYXJlYTogImFyY2hBIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiR29vZ2xlQWdlbnRzIiwgaWNvblNpemU6IDIwCgogICAgWm9uZSBAQ1JBR19QbGFuX1JldHJpZXZlIHsKICAgICAgdGl0bGU6ICJQbGFuICYgTXVsdGktSG9wIFJldHJpZXZlIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJTcGFya2xlcyIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfcGxhbm5lcjogIjEuIFF1ZXJ5UGxhbm5lciIgfCAiRGVjb21wb3NlcyBTdWItUXVlcmllcyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdvb2dsZUFnZW50cyIgfQogICAgICBbYWdlbnRfcmV0cmlldmVyOiAiMi4gSHlicmlkUmV0cmlldmVyIiB8ICIwLjcgQ29zaW5lICsgMC4zIEJNMjUiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJTZWFyY2giIH0KICAgIH0KCiAgICBab25lIEBDUkFHX0dyYWRlX1N5bnRoZXNpemUgewogICAgICB0aXRsZTogIkNSQUcgQ3JpdGljICYgU3ludGhlc2lzIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJWZXJ0ZXhBSSIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfZ3JhZGVyOiAiMy4gR3JhZGVyQ3JpdGljIiB8ICJGaWx0ZXJzIE5vaXNlIENodW5rcyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlNoaWVsZENoZWNrIiB9CiAgICAgIFthZ2VudF9zeW50aDogIjQuIFN5bnRoZXNpemVyIiB8ICJHcm91bmRlZCBTU0UgU3RyZWFtIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2VtaW5pIiB9CiAgICB9CiAgfQoKICBab25lIEBQYXRoX0JfQXJlbmFfQW5kX0V2YWwgewogICAgdGl0bGU6ICIyQi4gU3RhbmRhcmQgUkFHLCBNb2RlbCBBcmVuYSAmIEF1dG9yYXRlciBKdWRnZSIKICAgIHN1YnRpdGxlOiAiRmxhc2ggdnMuIFBybyBCZW5jaG1hcmsgJiBMTE0tYXMtYS1KdWRnZSAoLzUpIgogICAgc3R5bGU6IEBab25lQmx1ZSwgYXJlYTogImFyY2hCIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiVmVydGV4QUkiLCBpY29uU2l6ZTogMjAKCiAgICBab25lIEBNdWx0aU1vZGVsX0FyZW5hIHsKICAgICAgdGl0bGU6ICJNdWx0aS1Nb2RlbCBBcmVuYSIsIGxheW91dDogY29sdW1uLCBnYXA6IDE2LCBwYWRkaW5nOiAxNCwgYWxpZ246ICJjZW50ZXIiLAogICAgICBmaWxsOiAiI2ZmZmZmZiIsIHN0cm9rZTogIiM1YjliZjMiLCBjb3JuZXJSYWRpdXM6IDEwLCBpY29uOiAiR2VtaW5pIiwgaWNvblNpemU6IDE4CiAgICAgIFtnZW1pbmlfZmxhc2g6ICJHZW1pbmkgMy41IEZsYXNoIiB8ICJTdWItNDAwbXMgVFRGVCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdlbWluaSIgfQogICAgICBbZ2VtaW5pX3BybzogIkdlbWluaSAzLjEgUHJvIiB8ICJEZWVwIEZyb250aWVyIFJlYXNvbmluZyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlZlcnRleEFJIiB9CiAgICB9CgogICAgWm9uZSBAQXV0b3JhdGVyX0FuZF9GaW5vcHMgewogICAgICB0aXRsZTogIlF1YWxpdHkgJiBGaW5PcHMgQXVkaXQiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjNWI5YmYzIiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIkFjdGl2aXR5IiwgaWNvblNpemU6IDE4CiAgICAgIFt2ZXJ0ZXhfYXV0b3JhdGVyOiAiVmVydGV4IEF1dG9yYXRlciIgfCAiR3JvdW5kZWRuZXNzICYgUmVsZXZhbmNlIC81Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiU2VjdXJpdHlDb21tYW5kQ2VudGVyIiB9CiAgICAgIFtmaW5vcHNfYmFkZ2U6ICJSZWFsLVRpbWUgRmluT3BzIiB8ICJUb2tlbiBDb3N0ICh+JDAuMDAwMTEpIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRCaWxsaW5nIiB9CiAgICB9CiAgfQoKICBab25lIEBTaGFyZWRfRGF0YV9Gb3VuZGF0aW9uIHsKICAgIHRpdGxlOiAiMy4gU2hhcmVkIEh5YnJpZCBWZWN0b3IgSW5kZXggJiBDbG91ZCBTdG9yYWdlIFBlcnNpc3RlbmNlIChEaXJlY3QgVlBDIEVncmVzcykiCiAgICBzdWJ0aXRsZTogIkluLU1lbW9yeSA3NjgtZGltIENvc2luZSArIEJNMjUgSW5kZXggQmFja2VkIGJ5IENsb3VkIFN0b3JhZ2UgVUJMQSIKICAgIHN0eWxlOiBAWm9uZUdyZWVuLCBhcmVhOiAic3RvcmUiCiAgICBsYXlvdXQ6IHJvdywgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkdjcERhdGFiYXNlIiwgaWNvblNpemU6IDIwCiAgICBbZ2NzX2NvcnB1c19idWNrZXQ6ICJDbG91ZCBTdG9yYWdlIFVCTEEiIHwgIkNvcnB1cyAmIEpTT05MIEluZGV4Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogICAgW3ZlcnRleF9lbWJlZGRpbmdzOiAiVmVydGV4IEVtYmVkZGluZ3MiIHwgInRleHQtZW1iZWRkaW5nLTAwNCAoNzY4ZCkiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJBSVBsYXRmb3JtIiB9CiAgICBbaHlicmlkX2VuZ2luZTogIkdvIEh5YnJpZCBFbmdpbmUiIHwgIkNvc2luZSAoMC43KSArIEJNMjUgKDAuMykiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJDcHUiIH0KICAgIFtjbG91ZF9sb2dnaW5nX2V2YWw6ICJDbG91ZCBMb2dnaW5nIiB8ICJTU0UgJiBFdmFsIFRlbGVtZXRyeSJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkTG9nZ2luZyIgfQogIH0KfQoKW2FybW9yX2lhcF0gLT4gW2Nsb3VkX3J1bl9nb10geyBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiLCBzZXF1ZW5jZUJhZGdlOiAiMSIsIGJhZGdlRmlsbDogJEdjcEJsdWUsIGJhZGdlRm9udENvbG9yOiAkU3VyZmFjZVdoaXRlIH0KW2Nsb3VkX3J1bl9nb10gLT4gW2RvY19pbmdlc3Rvcl0geyBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2Nsb3VkX3J1bl9nb10gLT4gW2FnZW50X3BsYW5uZXJdIHsgbGFiZWw6ICJBZ2VudGljIFNTRSIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIyIiwgYmFkZ2VGaWxsOiAkR2NwQmx1ZSwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbYWdlbnRfcGxhbm5lcl0gLT4gW2FnZW50X3JldHJpZXZlcl0geyBsYWJlbDogIlN1Yi1RdWVyaWVzIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9yZXRyaWV2ZXJdIC0+IFthZ2VudF9ncmFkZXJdIHsgbGFiZWw6ICJUb3AtSyBDaHVua3MiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2FnZW50X2dyYWRlcl0gLT4gW2FnZW50X3N5bnRoXSB7IGxhYmVsOiAiVmVyaWZpZWQiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KClthZ2VudF9zeW50aF0gLT4gW2dlbWluaV9mbGFzaF0geyBsYWJlbDogIkNvbXBhcmUiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAxLjgsIGRhc2hlZDogdHJ1ZSwgY3VydmU6ICJzdGVwIiB9CltnZW1pbmlfZmxhc2hdIC0+IFtnZW1pbmlfcHJvXSB7IGxhYmVsOiAiQXJlbmEiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAxLjgsIGN1cnZlOiAic3RlcCIgfQpbZ2VtaW5pX2ZsYXNoXSAtPiBbdmVydGV4X2F1dG9yYXRlcl0geyBsYWJlbDogIkp1ZGdlIC81Iiwgc291cmNlQW5jaG9yOiAicmlnaHQiLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BHcmVlbiwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIzIiwgYmFkZ2VGaWxsOiAkR2NwR3JlZW4sIGJhZGdlRm9udENvbG9yOiAkU3VyZmFjZVdoaXRlIH0KW3ZlcnRleF9hdXRvcmF0ZXJdIC0+IFtmaW5vcHNfYmFkZ2VdIHsgbGFiZWw6ICJUb2tlbiBDb3N0Iiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BHcmVlbiwgc3Ryb2tlV2lkdGg6IDEuOCwgY3VydmU6ICJzdGVwIiB9CgpbZG9jX2luZ2VzdG9yXSAtPiBbZ2NzX2NvcnB1c19idWNrZXRdIHsgbGFiZWw6ICJQZXJzaXN0IEpTT05MIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfcmV0cmlldmVyXSAtPiBbdmVydGV4X2VtYmVkZGluZ3NdIHsgbGFiZWw6ICI3NjhkIFZlY3RvciIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfZ3JhZGVyXSAtPiBbaHlicmlkX2VuZ2luZV0geyBsYWJlbDogIkJNMjUgKyBDb3MiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2Zpbm9wc19iYWRnZV0gLT4gW2Nsb3VkX2xvZ2dpbmdfZXZhbF0geyBsYWJlbDogIkF1ZGl0IExvZ3MiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMS44LCBkYXNoZWQ6IHRydWUsIGN1cnZlOiAic3RlcCIgfQo=) · [**🖼️ PNG Haute-Résolution**](docs/diagrams/rag_comparison_crag_v3.png) · [`docs/diagrams/rag_comparison_crag_v3.dendrite`](docs/diagrams/rag_comparison_crag_v3.dendrite) · [`docs/diagrams/rag_comparison_crag_v3.drawio`](docs/diagrams/rag_comparison_crag_v3.drawio) · [`docs/architecture-gcpdraw.md`](docs/architecture-gcpdraw.md)

# RAG Comparison Demo

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Google Cloud](https://img.shields.io/badge/Google_Cloud-Vertex_AI-4285F4?style=flat&logo=google-cloud)](https://cloud.google.com/vertex-ai)
[![Gemini Models](https://img.shields.io/badge/Gemini-3.5_Flash_%7C_3.8_Flash_%7C_3.1_Pro-8E75B2?style=flat&logo=google-gemini)](https://ai.google.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Application web de démonstration technique comparant la recherche lexicale classique (mots-clés) et la génération augmentée par récupération (RAG groundé) avec la famille de modèles Google Gemini 3 sur Vertex AI.

L'application s'exécute comme un binaire Go autonome compilé sans dépendance externe (`0 dependencies`, `0 CVE`). Elle est déployée sur **Google Cloud Run** ou **Google Kubernetes Engine (GKE) Autopilot** au sein de l'architecture [GCP AI Foundation Blueprint](https://github.com/cloud-gtm/gcp-ai-foundation-blueprint).

---

## 🗺️ Architecture Globale (Dendrite v3.3 — Bento Matrix 16:9)

[![RAG Comparison — Hybrid RAG & 4-Agent CRAG Swarm Comparator](docs/diagrams/rag_comparison_crag_v3.png)](docs/diagrams/rag_comparison_crag_v3.png)

> 💡 **Astuce** : Cliquez sur l'image pour l'afficher en pleine résolution (Retina 2x), ou [**ouvrez-la dans Dendrite Studio ↗**](https://dendrite-758054785671.cr.gclb.goog/#code=dGhlbWU6ICJnY3AtYXJjaGl0ZWN0dXJlIgpyZW5kZXJPcmRlcjogbm9kZXMtZmlyc3QKZGlyZWN0aW9uOiByaWdodApzcGFjaW5nOiAyNgoKY29uc3QgR2NwQmx1ZSA9ICIjMWE3M2U4Igpjb25zdCBHY3BHcmVlbiA9ICIjMWU4ZTNlIgpjb25zdCBEYXJrU2xhdGUgPSAiIzIwMjEyNCIKY29uc3QgU3ViVGV4dCA9ICIjNWY2MzY4Igpjb25zdCBDYXJkQm9yZGVyID0gIiNkYWRjZTAiCmNvbnN0IFN1cmZhY2VXaGl0ZSA9ICIjZmZmZmZmIgoKU3R5bGUgQFByb2R1Y3RDYXJkIHsKICB3aWR0aDogMTk2LCBoZWlnaHQ6IDYwLAogIGZpbGw6ICRTdXJmYWNlV2hpdGUsIHN0cm9rZUNvbG9yOiAkQ2FyZEJvcmRlciwgc3Ryb2tlV2lkdGg6IDEsIGJvcmRlclJhZGl1czogOCwKICBmb250Q29sb3I6ICREYXJrU2xhdGUsIHN1YkZvbnRDb2xvcjogJFN1YlRleHQsCiAgZm9udFNpemU6IDEzLCBzdWJGb250U2l6ZTogMTAuNSwgbGFiZWxXZWlnaHQ6IGJvbGQsCiAgdGV4dEFsaWduOiAibGVmdCIsIHRleHRWQWxpZ246ICJtaWRkbGUiLAogIGljb25Qb3NpdGlvbjogImxlZnQiLCBpY29uU2l6ZTogMjYsIHBhZGRpbmc6IDEwLCBzaGFkb3c6IHRydWUKfQpTdHlsZSBAWm9uZUJsdWUgICB7IGZpbGw6ICIjZThmMGZlIiwgc3Ryb2tlOiAiIzViOWJmMyIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZUdyZWVuICB7IGZpbGw6ICIjZTZmNGVhIiwgc3Ryb2tlOiAiIzY4Yjg4ZSIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVBlYWNoICB7IGZpbGw6ICIjZmNlOGU2Iiwgc3Ryb2tlOiAiI2YyOGI4MiIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVB1cnBsZSB7IGZpbGw6ICIjZjNlOGZkIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGJvcmRlclJhZGl1czogMTIgfQoKWm9uZSBAUkFHX0NvbXBhcmlzb25fQXJjaGl0ZWN0dXJlIHsKICB0aXRsZTogIkdDUCBBSSBGb3VuZGF0aW9uIOKAlCBIeWJyaWQgUkFHICYgNC1BZ2VudCBDUkFHIFN3YXJtIENvbXBhcmF0b3IiCiAgc3VidGl0bGU6ICJDbG91ZCBSdW4gdjIgKEdvIDEuMjQgWmVyby1DVkUpIOKAoiBSZWFsLVRpbWUgU1NFIFN0cmVhbWluZyDigKIgVmVydGV4IEFJIEF1dG9yYXRlciAoZXVyb3BlLXdlc3QxKSIKICBpY29uOiAiR29vZ2xlQ2xvdWQiLCBpY29uU2l6ZTogMjQKICBsYXlvdXQ6IG1hdHJpeCwgZ2FwOiAyNCwgcGFkZGluZzogMjQsIGFsaWduOiAic3RyZXRjaCIKICBmaWxsOiAiI2Y4ZmFmZCIsIHN0cm9rZTogIiNkYWRjZTAiLCBjb3JuZXJSYWRpdXM6IDE0CiAgYXJlYXM6IFsKICAgICJ1cHN0cmVhbSAgYXJjaEEgIGFyY2hCIiwKICAgICJ1cHN0cmVhbSAgc3RvcmUgIHN0b3JlIgogIF0KCiAgWm9uZSBAVXBzdHJlYW1fRWRnZSB7CiAgICB0aXRsZTogIjEuIFplcm8tVHJ1c3QgSW5ncmVzcyIKICAgIHN1YnRpdGxlOiAiQ2xvdWQgUnVuIHYyIFNTRSIKICAgIHN0eWxlOiBAWm9uZVBlYWNoLCBhcmVhOiAidXBzdHJlYW0iCiAgICBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkNsb3VkQXJtb3IiLCBpY29uU2l6ZTogMjAKICAgIFthcm1vcl9pYXA6ICJDbG91ZCBBcm1vciAmIElBUCIgfCAiV0FGIEw3ICsgWmVyby1UcnVzdCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkQXJtb3IiIH0KICAgIFtjbG91ZF9ydW5fZ286ICJDbG91ZCBSdW4gdjIgKEdvKSIgfCAiU1NFIFJvdXRlciAmIFgtUmF5IFVJIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwQ2xvdWRSdW4iIH0KICAgIFtkb2NfaW5nZXN0b3I6ICJEb2N1bWVudCBJbmdlc3RvciIgfCAiUERGIENNYXAgLyBET0NYIC8gVFhUIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogIH0KCiAgWm9uZSBAUGF0aF9BX0FnZW50aWNfQ1JBRyB7CiAgICB0aXRsZTogIjJBLiBNb2RlIEFnZW50aWMgUkFHICg0LVN1YmFnZW50IENSQUcgU3dhcm0pIgogICAgc3VidGl0bGU6ICJNdWx0aS1Ib3AgRGVjb21wb3NpdGlvbiAmIFNlbGYtQ29ycmVjdGlvbiIKICAgIHN0eWxlOiBAWm9uZVB1cnBsZSwgYXJlYTogImFyY2hBIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiR29vZ2xlQWdlbnRzIiwgaWNvblNpemU6IDIwCgogICAgWm9uZSBAQ1JBR19QbGFuX1JldHJpZXZlIHsKICAgICAgdGl0bGU6ICJQbGFuICYgTXVsdGktSG9wIFJldHJpZXZlIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJTcGFya2xlcyIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfcGxhbm5lcjogIjEuIFF1ZXJ5UGxhbm5lciIgfCAiRGVjb21wb3NlcyBTdWItUXVlcmllcyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdvb2dsZUFnZW50cyIgfQogICAgICBbYWdlbnRfcmV0cmlldmVyOiAiMi4gSHlicmlkUmV0cmlldmVyIiB8ICIwLjcgQ29zaW5lICsgMC4zIEJNMjUiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJTZWFyY2giIH0KICAgIH0KCiAgICBab25lIEBDUkFHX0dyYWRlX1N5bnRoZXNpemUgewogICAgICB0aXRsZTogIkNSQUcgQ3JpdGljICYgU3ludGhlc2lzIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJWZXJ0ZXhBSSIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfZ3JhZGVyOiAiMy4gR3JhZGVyQ3JpdGljIiB8ICJGaWx0ZXJzIE5vaXNlIENodW5rcyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlNoaWVsZENoZWNrIiB9CiAgICAgIFthZ2VudF9zeW50aDogIjQuIFN5bnRoZXNpemVyIiB8ICJHcm91bmRlZCBTU0UgU3RyZWFtIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2VtaW5pIiB9CiAgICB9CiAgfQoKICBab25lIEBQYXRoX0JfQXJlbmFfQW5kX0V2YWwgewogICAgdGl0bGU6ICIyQi4gU3RhbmRhcmQgUkFHLCBNb2RlbCBBcmVuYSAmIEF1dG9yYXRlciBKdWRnZSIKICAgIHN1YnRpdGxlOiAiRmxhc2ggdnMuIFBybyBCZW5jaG1hcmsgJiBMTE0tYXMtYS1KdWRnZSAoLzUpIgogICAgc3R5bGU6IEBab25lQmx1ZSwgYXJlYTogImFyY2hCIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiVmVydGV4QUkiLCBpY29uU2l6ZTogMjAKCiAgICBab25lIEBNdWx0aU1vZGVsX0FyZW5hIHsKICAgICAgdGl0bGU6ICJNdWx0aS1Nb2RlbCBBcmVuYSIsIGxheW91dDogY29sdW1uLCBnYXA6IDE2LCBwYWRkaW5nOiAxNCwgYWxpZ246ICJjZW50ZXIiLAogICAgICBmaWxsOiAiI2ZmZmZmZiIsIHN0cm9rZTogIiM1YjliZjMiLCBjb3JuZXJSYWRpdXM6IDEwLCBpY29uOiAiR2VtaW5pIiwgaWNvblNpemU6IDE4CiAgICAgIFtnZW1pbmlfZmxhc2g6ICJHZW1pbmkgMy41IEZsYXNoIiB8ICJTdWItNDAwbXMgVFRGVCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdlbWluaSIgfQogICAgICBbZ2VtaW5pX3BybzogIkdlbWluaSAzLjEgUHJvIiB8ICJEZWVwIEZyb250aWVyIFJlYXNvbmluZyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlZlcnRleEFJIiB9CiAgICB9CgogICAgWm9uZSBAQXV0b3JhdGVyX0FuZF9GaW5vcHMgewogICAgICB0aXRsZTogIlF1YWxpdHkgJiBGaW5PcHMgQXVkaXQiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjNWI5YmYzIiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIkFjdGl2aXR5IiwgaWNvblNpemU6IDE4CiAgICAgIFt2ZXJ0ZXhfYXV0b3JhdGVyOiAiVmVydGV4IEF1dG9yYXRlciIgfCAiR3JvdW5kZWRuZXNzICYgUmVsZXZhbmNlIC81Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiU2VjdXJpdHlDb21tYW5kQ2VudGVyIiB9CiAgICAgIFtmaW5vcHNfYmFkZ2U6ICJSZWFsLVRpbWUgRmluT3BzIiB8ICJUb2tlbiBDb3N0ICh+JDAuMDAwMTEpIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRCaWxsaW5nIiB9CiAgICB9CiAgfQoKICBab25lIEBTaGFyZWRfRGF0YV9Gb3VuZGF0aW9uIHsKICAgIHRpdGxlOiAiMy4gU2hhcmVkIEh5YnJpZCBWZWN0b3IgSW5kZXggJiBDbG91ZCBTdG9yYWdlIFBlcnNpc3RlbmNlIChEaXJlY3QgVlBDIEVncmVzcykiCiAgICBzdWJ0aXRsZTogIkluLU1lbW9yeSA3NjgtZGltIENvc2luZSArIEJNMjUgSW5kZXggQmFja2VkIGJ5IENsb3VkIFN0b3JhZ2UgVUJMQSIKICAgIHN0eWxlOiBAWm9uZUdyZWVuLCBhcmVhOiAic3RvcmUiCiAgICBsYXlvdXQ6IHJvdywgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkdjcERhdGFiYXNlIiwgaWNvblNpemU6IDIwCiAgICBbZ2NzX2NvcnB1c19idWNrZXQ6ICJDbG91ZCBTdG9yYWdlIFVCTEEiIHwgIkNvcnB1cyAmIEpTT05MIEluZGV4Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogICAgW3ZlcnRleF9lbWJlZGRpbmdzOiAiVmVydGV4IEVtYmVkZGluZ3MiIHwgInRleHQtZW1iZWRkaW5nLTAwNCAoNzY4ZCkiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJBSVBsYXRmb3JtIiB9CiAgICBbaHlicmlkX2VuZ2luZTogIkdvIEh5YnJpZCBFbmdpbmUiIHwgIkNvc2luZSAoMC43KSArIEJNMjUgKDAuMykiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJDcHUiIH0KICAgIFtjbG91ZF9sb2dnaW5nX2V2YWw6ICJDbG91ZCBMb2dnaW5nIiB8ICJTU0UgJiBFdmFsIFRlbGVtZXRyeSJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkTG9nZ2luZyIgfQogIH0KfQoKW2FybW9yX2lhcF0gLT4gW2Nsb3VkX3J1bl9nb10geyBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiLCBzZXF1ZW5jZUJhZGdlOiAiMSIsIGJhZGdlRmlsbDogJEdjcEJsdWUsIGJhZGdlRm9udENvbG9yOiAkU3VyZmFjZVdoaXRlIH0KW2Nsb3VkX3J1bl9nb10gLT4gW2RvY19pbmdlc3Rvcl0geyBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2Nsb3VkX3J1bl9nb10gLT4gW2FnZW50X3BsYW5uZXJdIHsgbGFiZWw6ICJBZ2VudGljIFNTRSIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIyIiwgYmFkZ2VGaWxsOiAkR2NwQmx1ZSwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbYWdlbnRfcGxhbm5lcl0gLT4gW2FnZW50X3JldHJpZXZlcl0geyBsYWJlbDogIlN1Yi1RdWVyaWVzIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9yZXRyaWV2ZXJdIC0+IFthZ2VudF9ncmFkZXJdIHsgbGFiZWw6ICJUb3AtSyBDaHVua3MiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2FnZW50X2dyYWRlcl0gLT4gW2FnZW50X3N5bnRoXSB7IGxhYmVsOiAiVmVyaWZpZWQiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KClthZ2VudF9zeW50aF0gLT4gW2dlbWluaV9mbGFzaF0geyBsYWJlbDogIkNvbXBhcmUiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAxLjgsIGRhc2hlZDogdHJ1ZSwgY3VydmU6ICJzdGVwIiB9CltnZW1pbmlfZmxhc2hdIC0+IFtnZW1pbmlfcHJvXSB7IGxhYmVsOiAiQXJlbmEiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAxLjgsIGN1cnZlOiAic3RlcCIgfQpbZ2VtaW5pX2ZsYXNoXSAtPiBbdmVydGV4X2F1dG9yYXRlcl0geyBsYWJlbDogIkp1ZGdlIC81Iiwgc291cmNlQW5jaG9yOiAicmlnaHQiLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BHcmVlbiwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIzIiwgYmFkZ2VGaWxsOiAkR2NwR3JlZW4sIGJhZGdlRm9udENvbG9yOiAkU3VyZmFjZVdoaXRlIH0KW3ZlcnRleF9hdXRvcmF0ZXJdIC0+IFtmaW5vcHNfYmFkZ2VdIHsgbGFiZWw6ICJUb2tlbiBDb3N0Iiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BHcmVlbiwgc3Ryb2tlV2lkdGg6IDEuOCwgY3VydmU6ICJzdGVwIiB9CgpbZG9jX2luZ2VzdG9yXSAtPiBbZ2NzX2NvcnB1c19idWNrZXRdIHsgbGFiZWw6ICJQZXJzaXN0IEpTT05MIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfcmV0cmlldmVyXSAtPiBbdmVydGV4X2VtYmVkZGluZ3NdIHsgbGFiZWw6ICI3NjhkIFZlY3RvciIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfZ3JhZGVyXSAtPiBbaHlicmlkX2VuZ2luZV0geyBsYWJlbDogIkJNMjUgKyBDb3MiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2Zpbm9wc19iYWRnZV0gLT4gW2Nsb3VkX2xvZ2dpbmdfZXZhbF0geyBsYWJlbDogIkF1ZGl0IExvZ3MiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMS44LCBkYXNoZWQ6IHRydWUsIGN1cnZlOiAic3RlcCIgfQo=) pour zoomer et exporter en SVG / Draw.io.

---

## Fonctionnalités

### 1. Recherche Hybride (Sémantique & Lexicale)
- **Embeddings denses** : Vectorisation via le modèle Vertex AI `text-multilingual-embedding-002` (vecteurs de 768 dimensions).
- **Similarité cosinus en Go natif** : Calcul vectoriel sans base de données vectorielle externe ni bibliothèque CGO.
- **Score combiné** : Pondération $0{,}65 \times \text{Similarité}_{\text{cos}} + 0{,}35 \times \text{Score}_{\text{lex}}$ (BM25 normalisé).
- **Repli automatique** : Si l'API d'embeddings est indisponible ou hors quota, le moteur bascule de manière transparente sur la recherche lexicale avec normalisation linguistique (repliage d'accents, élisions, stop words).

### 2. Suite de Modèles Gemini 3 & Arena Multi-Modèles
- **Modèles supportés** :
  - `gemini-3.5-flash` (défaut) : Latence ultra-faible, débit élevé et efficience multimodale.
  - `gemini-3.8-flash` : Raisonnement de code et capacités agentiques avancées.
  - `gemini-3.1-pro-preview` : Raisonnement complexe et analyse approfondie d'architectures.
  - `gemini-3.5-flash-lite` : Cas d'usage haute fréquence à contrainte de coût minimale.
- **Sélecteur dynamique** : Changement de modèle à chaud via `/api/model/switch` sans redémarrage de conteneur.
- **Modes d'affichage (5 modes interactifs)** :
  - **RAG Simple** : Réponse unifiée avec passages sources et télémétrie.
  - **Split View** : Confrontation directe entre recherche lexicale brute et synthèse RAG.
  - **Arena (2 Modèles)** : Génération simultanée côte à côte sur la même requête avec deux modèles distincts.
  - **Triple Comparatif** : Écran 3 colonnes (Modèle A, Modèle B, Recherche Classique).
  - **🤖 Agentic RAG (Swarm 4 Sous-Agents CRAG)** : Pipeline multi-agents auto-correctif (`QueryPlanner` ➔ `HybridRetriever` ➔ `GraderCritic` avec boucle *Auto-Heal* ➔ `CitationSynthesizer`) accompagné d'un tiroir **Trace Live** affichant en direct chaque étape SSE (`event: agent_step`).
- **Télémétrie en temps réel** : Mesure du temps jusqu'au premier token (TTFT en ms), durée totale, et décompte exact des tokens via l'API Vertex AI.

### 3. Évaluation GenAI à la Demande (Autorater Vertex AI)
- **Déclenchement unitaire** : Bouton d'évaluation sous chaque réponse générée.
- **Métriques calculées** :
  - **Groundedness (Ancrage)** : Mesure l'alignement factuel de la réponse avec les passages documentaires fournis (note de 1 à 5).
  - **QA Relevance (Pertinence)** : Évalue la conformité de la réponse par rapport à la question posée (note de 1 à 5).
- **Rapports d'évaluation** : Modal interactif affichant les notes, les explications textuelles produites par l'Autorater Vertex AI, et les données de contexte évaluées.
- **Repli heuristique** : En cas d'indisponibilité du service Rapid Evaluation, une analyse locale par chevauchement n-grammes assure la continuité de service.

### 4. Persistance Cloud Storage & Multi-Instances
- **Index sérialisé** : Les documents ingérés et leurs embeddings sont persistés au format JSON compressé dans `gs://{GCS_RAG_BUCKET}/index/corpus.json`.
- **Scaling horizontal** : Prise en charge de plusieurs instances Cloud Run sans divergence d'état documentaire.
- **Démarrage instantané** : Restauration du snapshot GCS au boot de l'instance en moins d'une seconde.

---

## Architecture

Consultez le schéma d'architecture complet au format GCP Draw dans [docs/architecture-gcpdraw.md](docs/architecture-gcpdraw.md).

```
                            [ Navigateur Client ]
                                     │
                   HTTPS / Identity-Aware Proxy (IAP)
                                     │
                                     ▼
                     [ Google Cloud Load Balancer ]
                                     │
                                     ▼
                     [ Cloud Run : rag-comparison ]
                     ┌────────────────────────────┐
                     │ • Serveur HTTP Go natif    │
                     │ • Moteur Recherche Hybride │
                     │ • Swarm 4-Agents CRAG      │
                     │ • Static Assets (embed.FS) │
                     └──────┬──────────────┬──────┘
                            │              │
           Vectorisation &  │              │  Stockage Corpus
           Inférence LLM    │              │  & Snapshots
                            ▼              ▼
                    [ Vertex AI ]    [ Cloud Storage ]
                    • Gemini 3.5/3.8 • gs://{BUCKET}/index/
                    • Embeddings 002
                    • Rapid Eval
```

---

## Référence de l'API REST

Toutes les routes d'API sont servies par le binaire Go sur le port configuré (`PORT`, défaut `8080`).

| Méthode | Point de terminaison | Paramètres | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/documents` | Aucun | Liste l'ensemble des documents indexés avec leurs métadonnées (titre, taille, passages). |
| `POST` | `/api/documents/upload` | Multipart form (`files[]`) | Téléverse et indexe de nouveaux fichiers (PDF, texte, markdown) avec calcul d'embeddings. |
| `DELETE` | `/api/documents?id={id}` | `id` (identifiant document) | Supprime un document du corpus et met à jour l'index sur Cloud Storage. |
| `DELETE` | `/api/documents?all=true` | `all=true` | Purge l'intégralité du corpus documentaire en mémoire et sur Cloud Storage. |
| `GET` | `/api/search/classic` | `q={query}` | Exécute une recherche lexicale par mots-clés et retourne les extraits bruts. |
| `GET` | `/api/chat/stream` | `q={query}`, `model={model_id}`, `mode=agentic` *(opt.)* | Établit un flux SSE (`event: agent_step`, `retrieval`, `token`, `metrics`) pour streamer la réponse RAG standard ou le pipeline 4-agents CRAG. |
| `GET` | `/api/models` | Aucun | Liste les modèles Gemini disponibles, leurs caractéristiques et le modèle actif. |
| `POST` | `/api/model/switch` | Corps JSON `{"model": "id"}` | Modifie dynamiquement le modèle Gemini utilisé par défaut. |
| `POST` | `/api/evaluate` | Corps JSON `{"query", "prediction", "context"}` | Lance une évaluation d'ancrage et de pertinence via Vertex AI Rapid Evaluation. |
| `GET` | `/healthz` | Aucun | Sonde de vivacité et de disponibilité pour Cloud Run et Kubernetes. |

---

## Variables d'Environnement

L'application lit sa configuration depuis les variables d'environnement suivantes :

| Variable | Type | Valeur par défaut | Description |
| :--- | :--- | :--- | :--- |
| `GCP_PROJECT` (ou `PROJECT_ID`) | Chaîne | Découverte via métadonnées GCP | Identifiant du projet Google Cloud hébergeant Vertex AI. |
| `GCP_REGION` (ou `REGION`) | Chaîne | `europe-west1` | Région Google Cloud utilisée pour les appels d'API régionaux. |
| `GEMINI_MODEL` | Chaîne | `gemini-3.5-flash` | Modèle Gemini par défaut au démarrage du service. |
| `GCS_RAG_BUCKET` | Chaîne | Aucune (optionnel) | Nom du bucket GCS pour la persistance de l'index documentaire (`corpus.json`). |
| `PORT` | Entier | `8080` | Port d'écoute du serveur HTTP. |

---

## Démarrage Rapide

### Prérequis
- Go 1.25 ou supérieur installé localement.
- Google Cloud SDK (`gcloud`) authentifié avec accès au projet GCP cible.
- Droits IAM requis : `roles/aiplatform.user` et `roles/storage.objectAdmin`.

### Exécution Locale

```bash
# 1. Cloner le dépôt et se placer dans le répertoire source
cd src

# 2. Configurer les variables d'environnement GCP (GCP_PROJECT ou PROJECT_ID)
export GCP_PROJECT=$(gcloud config get-value project)
export GCP_REGION="europe-west1"
export GEMINI_MODEL="gemini-3.5-flash"
export GCS_RAG_BUCKET="${GCP_PROJECT}-rag-docs"

# 3. Lancer les tests unitaires
go test -v ./...

# 4. Démarrer le serveur applicatif
go run main.go
```

L'interface web est accessible sur `http://localhost:8080`.

---

## Déploiement

Consultez le [Guide de Déploiement Complet](docs/DEPLOYMENT_GUIDE.md) pour les instructions détaillées.

### Option 1 : Google Cloud Run (Serverless)

```bash
./deploy/cloudrun/deploy.sh
```

### Option 2 : GKE Autopilot (GCP AI Foundation Blueprint)

```bash
./scripts/deploy-to-blueprint.sh --blueprint-dir=../gcp-ai-foundation-blueprint
```

---

## Sécurité & Conformité

- **Gestion des identités** : Aucun secret ni clé de compte de service statique n'est embarqué dans l'image conteneur. L'authentification repose exclusivement sur **Workload Identity** (GKE) et le compte de service d'exécution Cloud Run via le serveur de métadonnées de l'instance (`http://metadata.google.internal`).
- **Contrôle d'accès réseau** : Déploiement recommandé avec l'option `--ingress=internal-and-cloud-load-balancing`, protégé en amont par Cloud Armor et Google Cloud Identity-Aware Proxy (IAP).
- **Zéro dépendance tierce** : Le code applicatif Go n'utilise aucun framework externe (`0 dépendances directes`, `0 vulnérabilités CVE` au scan conteneur).

---

## 🤖 Architecture Agentique Dual-Layer (Runtime CRAG Swarm & M1L1 Skills)

Ce dépôt implémente une architecture agentique à **deux niveaux complémentaires** :
- 🚀 **Couche 2 (Run-Time en Production)** : Un **Swarm de 4 Sous-Agents Corrective RAG (CRAG)** embarqué dans le binaire Go (`src/main.go`), déclenché en temps réel par l'utilisateur dans l'interface web.
- 🛠️ **Couche 1 (Build-Time en Ingénierie)** : **2 Sous-Agents spécialisés et 1 Skill M1L1** (`.agents/`), déclenchés dans l'IDE/CLI lors du développement et avant chaque `git commit`.

### 🔄 Diagramme de Séquence : Comment les 4 Sous-Agents CRAG entrent en action en direct

Lorsque l'utilisateur active le bouton **`🤖 Agentic RAG`** dans l'interface (`https://rag.hoffmannw.demo.altostrat.com`) et pose une question, voici l'enchaînement exact streamé en SSE (`event: agent_step`) vers le panneau **Trace Live** :

```mermaid
sequenceDiagram
    autonumber
    actor User as 👤 Utilisateur (UI Web)
    participant UI as 🖥️ Panneau Trace Live (SSE)
    participant A1 as 🧠 1. QueryPlannerAgent
    participant A2 as 🔍 2. HybridRetrieverAgent
    participant A3 as ⚖️ 3. GraderCriticAgent (CRAG)
    participant A4 as ✍️ 4. CitationSynthesizerAgent

    User->>UI: Clique "🤖 Agentic RAG" & envoie la question
    UI->>A1: GET /api/chat/stream?mode=agentic&q=...
    A1-->>UI: SSE agent_step (status: done, 3 sous-requêtes générées)
    A1->>A2: Transmet les sous-requêtes (lexicales + sémantiques)
    A2->>A2: Recherche parallèle Dense Cosine (768d) + BM25 (RRF k=60)
    A2-->>UI: SSE agent_step (status: done, N chunks dédupliqués)
    A2->>A3: Soumet les chunks candidats à l'audit factuel
    alt Couverture insuffisante (Score < 7/10) — Boucle Auto-Heal
        A3-->>UI: SSE agent_step (status: heal, réécriture + élargissement Top-K)
        A3->>A2: Relance HybridRetriever avec requête reformulée
        A2-->>A3: Nouveaux segments enrichis
    else Couverture validée (Score >= 7/10)
        A3-->>UI: SSE agent_step (status: done, pertinence validée)
    end
    A3->>A4: Transmet le contexte certifié
    A4-->>UI: Stream SSE token par token + citations [Doc, Chunk #X]
```

### 📊 Matrice Récapitulative : Où et Comment chaque Agent intervient

| Agent / Skill | Couche | Où vit-il ? | Comment / Quand entre-t-il en action ? | Rôle & Valeur ajoutée |
| :--- | :--- | :--- | :--- | :--- |
| **`QueryPlannerAgent`** | **Couche 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Étape 1** dès l'envoi d'une question en mode `🤖 Agentic RAG`. | Décompose une question complexe ou multi-critères en sous-requêtes ciblées (mots-clés lexicaux + concepts sémantiques). |
| **`HybridRetrieverAgent`** | **Couche 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Étape 2** après la planification (et lors d'un *Auto-Heal*). | Exécute la recherche hybride parallèle **Dense Cosine + Sparse BM25 via Reciprocal Rank Fusion ($k=60$)** et déduplique les passages. |
| **`GraderCriticAgent`** | **Couche 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Étape 3** avant toute génération LLM. | Évalue la pertinence factuelle des passages extraits (`/10`). Déclenche automatiquement une boucle **Auto-Heal (Query Rewrite)** si le contexte est incomplet. |
| **`CitationSynthesizerAgent`** | **Couche 2** *(Run-Time)* | `src/main.go` (`runAgenticRAGPipeline`) | **Étape 4** une fois les passages validés par le `GraderCriticAgent`. | Génère la synthèse en streaming SSE avec citations vérifiées `[Doc: <Titre>, Chunk #X]` et lance l'évaluation Autorater Vertex AI. |
| **[`rag-eval-scientist`](.agents/agents/rag-eval-scientist.md)** | **Couche 1** *(Build-Time)* | `.agents/agents/rag-eval-scientist.md` | Dans **Jetski / Antigravity / Gemini CLI** lors du tuning des poids RRF/BM25 ou des prompts d'évaluation. | Audite les hyperparamètres **RRF ($k=60$)**, BM25 ($k_1=1.2, b=0.75$) et la fidélité **LLM-as-a-Judge** (`Groundedness` / `QA Relevance`). |
| **[`go-concurrency-reviewer`](.agents/agents/go-concurrency-reviewer.md)** | **Couche 1** *(Build-Time)* | `.agents/agents/go-concurrency-reviewer.md` | Dans **Jetski / Antigravity / Gemini CLI** avant de commiter du code Go (`src/main.go`). | Vérifie l'absence de deadlock `sync.RWMutex`, la fermeture propre des flux SSE (`r.Context().Done()`) et la synchro GCS asynchrone. |
| **[`rag-benchmark-and-ci`](.agents/skills/rag-benchmark-and-ci/SKILL.md)** | **Couche 1** *(Gatekeeper)* | `.agents/skills/rag-benchmark-and-ci/scripts/verify.sh` | Exécuté dans le terminal avant chaque `git commit` ou déploiement Cloud Run. | Vérifie `go vet ./...`, `go test -v -race ./...`, la règle **Zero External Dependencies** (`src/go.mod`) et la présence des 4 agents CRAG. |

### 🎬 Scénario de Démo Client en 3 Minutes (Playbook CE)

1. **Étape 1 — Déclencher le Swarm CRAG en direct dans le navigateur (Run-Time)** :
   - Ouvrez **[RAG Comparison Demo](https://rag.hoffmannw.demo.altostrat.com)** et cliquez sur le bouton **`🤖 Agentic RAG`** dans la barre supérieure.
   - Copiez-collez une question multi-critères dans la barre de recherche :
     > `"Compare l'architecture de sécurité Zero-Trust (IAP, Cloud Armor) avec la stratégie de sauvegarde WORM et explique comment le RAG hybride évite les hallucinations."`
   - **Ce qu'il faut montrer à l'écran** :
     - Le panneau **Trace Live** s'ouvre automatiquement au-dessus de la réponse et affiche en temps réel les 4 cartes d'agents (`QueryPlanner` ➔ `HybridRetriever` ➔ `GraderCritic` ➔ `CitationSynthesizer`) avec leur latence en millisecondes et les sous-requêtes générées.
     - Cliquez ensuite sur **`⭐ Évaluer (Vertex AI)`** sous la réponse pour afficher le score **Groundedness (`/5`)** et **QA Relevance (`/5`)**.

2. **Étape 2 — Déclencher un Agent en ligne de commande via `curl` (SSE Stream)** :
   ```bash
   curl -N "https://rag.hoffmannw.demo.altostrat.com/api/chat/stream?mode=agentic&q=Architecture+Cloud+Armor+et+RRF"
   ```
   *(Affiche en direct les événements `event: agent_step` JSON suivis des tokens streamés).*

3. **Étape 3 — Démontrer les Agents d'Ingénierie & le Gatekeeper M1L1 (Build-Time)** :
   - Dans **Jetski / Antigravity / Gemini CLI**, copiez-collez :
     > `"Invoque rag-eval-scientist pour auditer la formule Reciprocal Rank Fusion (k=60) et les seuils du GraderCriticAgent dans src/main.go."`
   - Puis lancez le script gatekeeper M1L1 :
     ```bash
     ./.agents/skills/rag-benchmark-and-ci/scripts/verify.sh
     ```

---

## Licence

Ce projet est distribué sous licence Apache 2.0. Consultez le fichier [LICENSE](LICENSE) pour plus d'informations.



