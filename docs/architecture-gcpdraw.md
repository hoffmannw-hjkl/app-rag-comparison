# 🏗️ Schémas d'Architecture Officiels (Dendrite v3.3 Bento Matrix & GCP Draw)

Ce document centralise les **schémas d'architecture as-code** de **RAG Comparison Platform — 4-Engine Benchmark, CRAG Self-Correction & LLM-as-a-Judge** selon deux standards complémentaires :
1. **Dendrite v3.3 (Executive Standard 16:9 — Bento Matrix)** : Schéma haute-densité validé sans collision géométrique, exportable en **SVG interactif** et **Draw.io / Lucidchart (`.drawio`)**.
2. **GCP Draw (`go/gcpdraw`)** : Spécification déclarative rapide pour l'outil interne Google Cloud Draw.

---

## 🚀 Accès Rapide & Fichiers Sources Versionnés

| Format / Outil | Fichier / Lien Direct | Usage Recommandé |
| :--- | :--- | :--- |
| **🎨 Ouvrir dans Dendrite Studio (1-Click)** | [**Launch in Dendrite Studio ↗**](https://dendrite-758054785671.cr.gclb.goog/#code=dGhlbWU6ICJnY3AtYXJjaGl0ZWN0dXJlIgpyZW5kZXJPcmRlcjogbm9kZXMtZmlyc3QKZGlyZWN0aW9uOiByaWdodApzcGFjaW5nOiAyNgoKY29uc3QgR2NwQmx1ZSA9ICIjMWE3M2U4Igpjb25zdCBHY3BHcmVlbiA9ICIjMWU4ZTNlIgpjb25zdCBEYXJrU2xhdGUgPSAiIzIwMjEyNCIKY29uc3QgU3ViVGV4dCA9ICIjNWY2MzY4Igpjb25zdCBDYXJkQm9yZGVyID0gIiNkYWRjZTAiCmNvbnN0IFN1cmZhY2VXaGl0ZSA9ICIjZmZmZmZmIgoKU3R5bGUgQFByb2R1Y3RDYXJkIHsKICB3aWR0aDogMTk2LCBoZWlnaHQ6IDYwLAogIGZpbGw6ICRTdXJmYWNlV2hpdGUsIHN0cm9rZUNvbG9yOiAkQ2FyZEJvcmRlciwgc3Ryb2tlV2lkdGg6IDEsIGJvcmRlclJhZGl1czogOCwKICBmb250Q29sb3I6ICREYXJrU2xhdGUsIHN1YkZvbnRDb2xvcjogJFN1YlRleHQsCiAgZm9udFNpemU6IDEzLCBzdWJGb250U2l6ZTogMTAuNSwgbGFiZWxXZWlnaHQ6IGJvbGQsCiAgdGV4dEFsaWduOiAibGVmdCIsIHRleHRWQWxpZ246ICJtaWRkbGUiLAogIGljb25Qb3NpdGlvbjogImxlZnQiLCBpY29uU2l6ZTogMjYsIHBhZGRpbmc6IDEwLCBzaGFkb3c6IHRydWUKfQpTdHlsZSBAWm9uZUJsdWUgICB7IGZpbGw6ICIjZThmMGZlIiwgc3Ryb2tlOiAiIzViOWJmMyIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZUdyZWVuICB7IGZpbGw6ICIjZTZmNGVhIiwgc3Ryb2tlOiAiIzY4Yjg4ZSIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVBlYWNoICB7IGZpbGw6ICIjZmNlOGU2Iiwgc3Ryb2tlOiAiI2YyOGI4MiIsIGJvcmRlclJhZGl1czogMTIgfQpTdHlsZSBAWm9uZVB1cnBsZSB7IGZpbGw6ICIjZjNlOGZkIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGJvcmRlclJhZGl1czogMTIgfQoKWm9uZSBAUkFHX0NvbXBhcmlzb25fQXJjaGl0ZWN0dXJlIHsKICB0aXRsZTogIkdDUCBBSSBGb3VuZGF0aW9uIOKAlCBIeWJyaWQgUkFHICYgNC1BZ2VudCBDUkFHIFN3YXJtIENvbXBhcmF0b3IiCiAgc3VidGl0bGU6ICJDbG91ZCBSdW4gdjIgKEdvIDEuMjQgWmVyby1DVkUpIOKAoiBSZWFsLVRpbWUgU1NFIFN0cmVhbWluZyDigKIgVmVydGV4IEFJIEF1dG9yYXRlciAoZXVyb3BlLXdlc3QxKSIKICBpY29uOiAiR29vZ2xlQ2xvdWQiLCBpY29uU2l6ZTogMjQKICBsYXlvdXQ6IG1hdHJpeCwgZ2FwOiAyNCwgcGFkZGluZzogMjQsIGFsaWduOiAic3RyZXRjaCIKICBmaWxsOiAiI2Y4ZmFmZCIsIHN0cm9rZTogIiNkYWRjZTAiLCBjb3JuZXJSYWRpdXM6IDE0CiAgYXJlYXM6IFsKICAgICJ1cHN0cmVhbSAgYXJjaEEgIGFyY2hCIiwKICAgICJ1cHN0cmVhbSAgc3RvcmUgIHN0b3JlIgogIF0KCiAgWm9uZSBAVXBzdHJlYW1fRWRnZSB7CiAgICB0aXRsZTogIjEuIFplcm8tVHJ1c3QgSW5ncmVzcyIKICAgIHN1YnRpdGxlOiAiQ2xvdWQgUnVuIHYyIFNTRSIKICAgIHN0eWxlOiBAWm9uZVBlYWNoLCBhcmVhOiAidXBzdHJlYW0iCiAgICBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkNsb3VkQXJtb3IiLCBpY29uU2l6ZTogMjAKICAgIFthcm1vcl9pYXA6ICJDbG91ZCBBcm1vciAmIElBUCIgfCAiV0FGIEw3ICsgWmVyby1UcnVzdCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkNsb3VkQXJtb3IiIH0KICAgIFtjbG91ZF9ydW5fZ286ICJDbG91ZCBSdW4gdjIgKEdvKSIgfCAiU1NFIFJvdXRlciAmIFgtUmF5IFVJIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwQ2xvdWRSdW4iIH0KICAgIFtkb2NfaW5nZXN0b3I6ICJEb2N1bWVudCBJbmdlc3RvciIgfCAiUERGIENNYXAgLyBET0NYIC8gVFhUIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogIH0KCiAgWm9uZSBAUGF0aF9BX0FnZW50aWNfQ1JBRyB7CiAgICB0aXRsZTogIjJBLiBNb2RlIEFnZW50aWMgUkFHICg0LVN1YmFnZW50IENSQUcgU3dhcm0pIgogICAgc3VidGl0bGU6ICJNdWx0aS1Ib3AgRGVjb21wb3NpdGlvbiAmIFNlbGYtQ29ycmVjdGlvbiIKICAgIHN0eWxlOiBAWm9uZVB1cnBsZSwgYXJlYTogImFyY2hBIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiR29vZ2xlQWdlbnRzIiwgaWNvblNpemU6IDIwCgogICAgWm9uZSBAQ1JBR19QbGFuX1JldHJpZXZlIHsKICAgICAgdGl0bGU6ICJQbGFuICYgTXVsdGktSG9wIFJldHJpZXZlIiwgbGF5b3V0OiBjb2x1bW4sIGdhcDogMTYsIHBhZGRpbmc6IDE0LCBhbGlnbjogImNlbnRlciIsCiAgICAgIGZpbGw6ICIjZmZmZmZmIiwgc3Ryb2tlOiAiI2ExNDJmNCIsIGNvcm5lclJhZGl1czogMTAsIGljb246ICJTcGFya2xlcyIsIGljb25TaXplOiAxOAogICAgICBbYWdlbnRfcGxhbm5lcjogIjEuIFF1ZXJ5UGxhbm5lckFnZW50IiB8ICJEZWNvbXBvc2VzIFN1Yi1RdWVyaWVzIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR29vZ2xlQWdlbnRzIiB9CiAgICAgIFthZ2VudF9yZXRyaWV2ZXI6ICIyLiBIeWJyaWRSZXRyaWV2ZXIiIHwgIjAuNyBDb3NpbmUgKyAwLjMgQk0yNSJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlNlYXJjaCIgfQogICAgfQoKICAgIFpvbmUgQENSQUdfR3JhZGVfU3ludGhlc2l6ZSB7CiAgICAgIHRpdGxlOiAiQ1JBRyBDcml0aWMgJiBTeW50aGVzaXMiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjYTE0MmY0IiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIlZlcnRleEFJIiwgaWNvblNpemU6IDE4CiAgICAgIFthZ2VudF9ncmFkZXI6ICIzLiBHcmFkZXJDcml0aWNBZ2VudCIgfCAiRmlsdGVycyBOb2lzZSBDaHVua3MiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJTaGllbGRDaGVjayIgfQogICAgICBbYWdlbnRfc3ludGg6ICI0LiBDaXRhdGlvblN5bnRoZXNpemVyIiB8ICJHcm91bmRlZCBTU0UgU3RyZWFtIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2VtaW5pIiB9CiAgICB9CiAgfQoKICBab25lIEBQYXRoX0JfQXJlbmFfQW5kX0V2YWwgewogICAgdGl0bGU6ICIyQi4gU3RhbmRhcmQgUkFHLCBNb2RlbCBBcmVuYSAmIEF1dG9yYXRlciBKdWRnZSIKICAgIHN1YnRpdGxlOiAiRmxhc2ggdnMuIFBybyBCZW5jaG1hcmsgJiBMTE0tYXMtYS1KdWRnZSAoLzUpIgogICAgc3R5bGU6IEBab25lQmx1ZSwgYXJlYTogImFyY2hCIgogICAgbGF5b3V0OiByb3csIGdhcDogMjIsIHBhZGRpbmc6IDE4LCBhbGlnbjogInN0cmV0Y2giLCBpY29uOiAiVmVydGV4QUkiLCBpY29uU2l6ZTogMjAKCiAgICBab25lIEBNdWx0aU1vZGVsX0FyZW5hIHsKICAgICAgdGl0bGU6ICJNdWx0aS1Nb2RlbCBBcmVuYSIsIGxheW91dDogY29sdW1uLCBnYXA6IDE2LCBwYWRkaW5nOiAxNCwgYWxpZ246ICJjZW50ZXIiLAogICAgICBmaWxsOiAiI2ZmZmZmZiIsIHN0cm9rZTogIiM1YjliZjMiLCBjb3JuZXJSYWRpdXM6IDEwLCBpY29uOiAiR2VtaW5pIiwgaWNvblNpemU6IDE4CiAgICAgIFtnZW1pbmlfZmxhc2g6ICJHZW1pbmkgMy41IC8gMy44IEZsYXNoIiB8ICJTdWItNDAwbXMgVFRGVCJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIkdlbWluaSIgfQogICAgICBbZ2VtaW5pX3BybzogIkdlbWluaSAzLjEgUHJvIiB8ICJEZWVwIEZyb250aWVyIFJlYXNvbmluZyJdIHsgc3R5bGU6IEBQcm9kdWN0Q2FyZCwgaWNvbjogIlZlcnRleEFJIiB9CiAgICB9CgogICAgWm9uZSBAQXV0b3JhdGVyX0FuZF9GaW5vcHMgewogICAgICB0aXRsZTogIlF1YWxpdHkgJiBGaW5PcHMgQXVkaXQiLCBsYXlvdXQ6IGNvbHVtbiwgZ2FwOiAxNiwgcGFkZGluZzogMTQsIGFsaWduOiAiY2VudGVyIiwKICAgICAgZmlsbDogIiNmZmZmZmYiLCBzdHJva2U6ICIjNWI5YmYzIiwgY29ybmVyUmFkaXVzOiAxMCwgaWNvbjogIkFjdGl2aXR5IiwgaWNvblNpemU6IDE4CiAgICAgIFt2ZXJ0ZXhfYXV0b3JhdGVyOiAiVmVydGV4IEFJIEF1dG9yYXRlciIgfCAiR3JvdW5kZWRuZXNzICYgUmVsZXZhbmNlIC81Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiU2VjdXJpdHlDb21tYW5kQ2VudGVyIiB9CiAgICAgIFtmaW5vcHNfYmFkZ2U6ICJSZWFsLVRpbWUgRmluT3BzIiB8ICJUb2tlbiBDb3N0ICh+JDAuMDAwMTEpIl0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRCaWxsaW5nIiB9CiAgICB9CiAgfQoKICBab25lIEBTaGFyZWRfRGF0YV9Gb3VuZGF0aW9uIHsKICAgIHRpdGxlOiAiMy4gU2hhcmVkIEh5YnJpZCBWZWN0b3IgSW5kZXggJiBDbG91ZCBTdG9yYWdlIFBlcnNpc3RlbmNlIChEaXJlY3QgVlBDIEVncmVzcykiCiAgICBzdWJ0aXRsZTogIkluLU1lbW9yeSA3NjgtZGltIENvc2luZSArIEJNMjUgSW5kZXggQmFja2VkIGJ5IENsb3VkIFN0b3JhZ2UgVUJMQSIKICAgIHN0eWxlOiBAWm9uZUdyZWVuLCBhcmVhOiAic3RvcmUiCiAgICBsYXlvdXQ6IHJvdywgZ2FwOiAyMCwgcGFkZGluZzogMTgsIGFsaWduOiAiY2VudGVyIiwgaWNvbjogIkdjcERhdGFiYXNlIiwgaWNvblNpemU6IDIwCiAgICBbZ2NzX2NvcnB1c19idWNrZXQ6ICJDbG91ZCBTdG9yYWdlIFVCTEEiIHwgIkNvcnB1cyAmIEpTT05MIEluZGV4Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiR2NwU3RvcmFnZUJ1Y2tldCIgfQogICAgW3ZlcnRleF9lbWJlZGRpbmdzOiAiVmVydGV4IEVtYmVkZGluZ3MiIHwgInRleHQtZW1iZWRkaW5nLTAwNCAoNzY4ZCkiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJBSVBsYXRmb3JtIiB9CiAgICBbaHlicmlkX2VuZ2luZTogIkdvIEh5YnJpZCBFbmdpbmUiIHwgIkNvc2luZSAoMC43KSArIEJNMjUgKDAuMykiXSB7IHN0eWxlOiBAUHJvZHVjdENhcmQsIGljb246ICJDcHUiIH0KICAgIFtjbG91ZF9sb2dnaW5nX2V2YWw6ICJDbG91ZCBMb2dnaW5nICYgVHJhY2UiIHwgIlNTRSAmIEV2YWwgVGVsZW1ldHJ5Il0geyBzdHlsZTogQFByb2R1Y3RDYXJkLCBpY29uOiAiQ2xvdWRMb2dnaW5nIiB9CiAgfQp9CgpbYXJtb3JfaWFwXSAtPiBbY2xvdWRfcnVuX2dvXSB7IHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIsIHNlcXVlbmNlQmFkZ2U6ICIxIiwgYmFkZ2VGaWxsOiAkR2NwQmx1ZSwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbY2xvdWRfcnVuX2dvXSAtPiBbZG9jX2luZ2VzdG9yXSB7IHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbY2xvdWRfcnVuX2dvXSAtPiBbYWdlbnRfcGxhbm5lcl0geyBsYWJlbDogIkFnZW50aWMgU1NFIiwgc291cmNlQW5jaG9yOiAicmlnaHQiLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiwgc2VxdWVuY2VCYWRnZTogIjIiLCBiYWRnZUZpbGw6ICRHY3BCbHVlLCBiYWRnZUZvbnRDb2xvcjogJFN1cmZhY2VXaGl0ZSB9ClthZ2VudF9wbGFubmVyXSAtPiBbYWdlbnRfcmV0cmlldmVyXSB7IGxhYmVsOiAiU3ViLVF1ZXJpZXMiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEJsdWUsIHN0cm9rZVdpZHRoOiAyLCBjdXJ2ZTogInN0ZXAiIH0KW2FnZW50X3JldHJpZXZlcl0gLT4gW2FnZW50X2dyYWRlcl0geyBsYWJlbDogIlRvcC1LIENodW5rcyIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbYWdlbnRfZ3JhZGVyXSAtPiBbYWdlbnRfc3ludGhdIHsgbGFiZWw6ICJWZXJpZmllZCIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQoKW2FnZW50X3N5bnRoXSAtPiBbZ2VtaW5pX2ZsYXNoXSB7IGxhYmVsOiAiQ29tcGFyZSIsIHNvdXJjZUFuY2hvcjogInJpZ2h0IiwgdGFyZ2V0QW5jaG9yOiAibGVmdCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDEuOCwgZGFzaGVkOiB0cnVlLCBjdXJ2ZTogInN0ZXAiIH0KW2dlbWluaV9mbGFzaF0gLT4gW2dlbWluaV9wcm9dIHsgbGFiZWw6ICJBcmVuYSIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDEuOCwgY3VydmU6ICJzdGVwIiB9CltnZW1pbmlfZmxhc2hdIC0+IFt2ZXJ0ZXhfYXV0b3JhdGVyXSB7IGxhYmVsOiAiSnVkZ2UgLzUiLCBzb3VyY2VBbmNob3I6ICJyaWdodCIsIHRhcmdldEFuY2hvcjogImxlZnQiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiwgc2VxdWVuY2VCYWRnZTogIjMiLCBiYWRnZUZpbGw6ICRHY3BHcmVlbiwgYmFkZ2VGb250Q29sb3I6ICRTdXJmYWNlV2hpdGUgfQpbdmVydGV4X2F1dG9yYXRlcl0gLT4gW2Zpbm9wc19iYWRnZV0geyBsYWJlbDogIlRva2VuIENvc3QiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJ0b3AiLCBjb2xvcjogJEdjcEdyZWVuLCBzdHJva2VXaWR0aDogMS44LCBjdXJ2ZTogInN0ZXAiIH0KCltkb2NfaW5nZXN0b3JdIC0+IFtnY3NfY29ycHVzX2J1Y2tldF0geyBsYWJlbDogIlBlcnNpc3QgSlNPTkwiLCBzb3VyY2VBbmNob3I6ICJib3R0b20iLCB0YXJnZXRBbmNob3I6ICJsZWZ0IiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9yZXRyaWV2ZXJdIC0+IFt2ZXJ0ZXhfZW1iZWRkaW5nc10geyBsYWJlbDogIjc2OGQgVmVjdG9yIiwgc291cmNlQW5jaG9yOiAiYm90dG9tIiwgdGFyZ2V0QW5jaG9yOiAidG9wIiwgY29sb3I6ICRHY3BCbHVlLCBzdHJva2VXaWR0aDogMiwgY3VydmU6ICJzdGVwIiB9ClthZ2VudF9ncmFkZXJdIC0+IFtoeWJyaWRfZW5naW5lXSB7IGxhYmVsOiAiQk0yNSArIENvcyIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwQmx1ZSwgc3Ryb2tlV2lkdGg6IDIsIGN1cnZlOiAic3RlcCIgfQpbZmlub3BzX2JhZGdlXSAtPiBbY2xvdWRfbG9nZ2luZ19ldmFsXSB7IGxhYmVsOiAiQXVkaXQgTG9ncyIsIHNvdXJjZUFuY2hvcjogImJvdHRvbSIsIHRhcmdldEFuY2hvcjogInRvcCIsIGNvbG9yOiAkR2NwR3JlZW4sIHN0cm9rZVdpZHRoOiAxLjgsIGRhc2hlZDogdHJ1ZSwgY3VydmU6ICJzdGVwIiB9Cg==) | Visualisation interactive plein écran, zoom vectoriel, export PNG/SVG/Draw.io en 1 clic |
| **📐 Source Dendrite v3.3 (`.dendrite`)** | [`docs/diagrams/rag_comparison_crag_v3.dendrite`](diagrams/rag_comparison_crag_v3.dendrite) | Code source déclaratif YAML/DSL Dendrite v3.3 (grille Bento 16:9, icônes GCP officielles) |
| **🧩 Export Draw.io / Diagrams.net (`.drawio`)** | [`docs/diagrams/rag_comparison_crag_v3.drawio`](diagrams/rag_comparison_crag_v3.drawio) | Fichier XML natif importable dans **Draw.io**, **Lucidchart** ou **Google Slides** |

---

## 🖼️ Aperçu Visuel (Dendrite v3.3 — Bento Matrix 16:9)

[![RAG Comparison — Hybrid RAG & 4-Agent CRAG Swarm Comparator](diagrams/rag_comparison_crag_v3.png)](diagrams/rag_comparison_crag_v3.png)

---

## 1. Spécification Dendrite v3.3 (`docs/diagrams/rag_comparison_crag_v3.dendrite`)

> **Sous-titre exécutif** : *Go 1.24 Serverless Engine on Cloud Run Gen2, Vertex AI Gemini 2.5 Flash, text-embedding-004 (768d), RRF k=60, CRAG & FinOps Telemetry*

Pour re-compiler ce schéma en local ou vérifier les contraintes géométriques (0 overlap, 0 clipping) :
```bash
node /google/src/files/head/depot/google3/cloud/professional_services/agents/skills/dendrite/scripts/render_dendrite.mjs \
  docs/diagrams/rag_comparison_crag_v3.dendrite
```

```yaml
theme: "gcp-architecture"
renderOrder: nodes-first
direction: right
spacing: 26

const GcpBlue = "#1a73e8"
const GcpGreen = "#1e8e3e"
const DarkSlate = "#202124"
const SubText = "#5f6368"
const CardBorder = "#dadce0"
const SurfaceWhite = "#ffffff"

Style @ProductCard {
  width: 196, height: 60,
  fill: $SurfaceWhite, strokeColor: $CardBorder, strokeWidth: 1, borderRadius: 8,
  fontColor: $DarkSlate, subFontColor: $SubText,
  fontSize: 13, subFontSize: 10.5, labelWeight: bold,
  textAlign: "left", textVAlign: "middle",
  iconPosition: "left", iconSize: 26, padding: 10, shadow: true
}
Style @ZoneBlue   { fill: "#e8f0fe", stroke: "#5b9bf3", borderRadius: 12 }
Style @ZoneGreen  { fill: "#e6f4ea", stroke: "#68b88e", borderRadius: 12 }
Style @ZonePeach  { fill: "#fce8e6", stroke: "#f28b82", borderRadius: 12 }
Style @ZonePurple { fill: "#f3e8fd", stroke: "#a142f4", borderRadius: 12 }

Zone @RAG_Comparison_Architecture {
  title: "GCP AI Foundation — Hybrid RAG & 4-Agent CRAG Swarm Comparator"
  subtitle: "Cloud Run v2 (Go 1.24 Zero-CVE) • Real-Time SSE Streaming • Vertex AI Autorater (europe-west1)"
  icon: "GoogleCloud", iconSize: 24
  layout: matrix, gap: 24, padding: 24, align: "stretch"
  fill: "#f8fafd", stroke: "#dadce0", cornerRadius: 14
  areas: [
    "upstream  archA  archB",
    "upstream  store  store"
  ]

  Zone @Upstream_Edge {
    title: "1. Zero-Trust Ingress"
    subtitle: "Cloud Run v2 SSE"
    style: @ZonePeach, area: "upstream"
    layout: column, gap: 20, padding: 18, align: "center", icon: "CloudArmor", iconSize: 20
    [armor_iap: "Cloud Armor & IAP" | "WAF L7 + Zero-Trust"] { style: @ProductCard, icon: "CloudArmor" }
    [cloud_run_go: "Cloud Run v2 (Go)" | "SSE Router & X-Ray UI"] { style: @ProductCard, icon: "GcpCloudRun" }
    [doc_ingestor: "Document Ingestor" | "PDF CMap / DOCX / TXT"] { style: @ProductCard, icon: "GcpStorageBucket" }
  }

  Zone @Path_A_Agentic_CRAG {
    title: "2A. Mode Agentic RAG (4-Subagent CRAG Swarm)"
    subtitle: "Multi-Hop Decomposition & Self-Correction"
    style: @ZonePurple, area: "archA"
    layout: row, gap: 22, padding: 18, align: "stretch", icon: "GoogleAgents", iconSize: 20

    Zone @CRAG_Plan_Retrieve {
      title: "Plan & Multi-Hop Retrieve", layout: column, gap: 16, padding: 14, align: "center",
      fill: "#ffffff", stroke: "#a142f4", cornerRadius: 10, icon: "Sparkles", iconSize: 18
      [agent_planner: "1. QueryPlannerAgent" | "Decomposes Sub-Queries"] { style: @ProductCard, icon: "GoogleAgents" }
      [agent_retriever: "2. HybridRetriever" | "0.7 Cosine + 0.3 BM25"] { style: @ProductCard, icon: "Search" }
    }

    Zone @CRAG_Grade_Synthesize {
      title: "CRAG Critic & Synthesis", layout: column, gap: 16, padding: 14, align: "center",
      fill: "#ffffff", stroke: "#a142f4", cornerRadius: 10, icon: "VertexAI", iconSize: 18
      [agent_grader: "3. GraderCriticAgent" | "Filters Noise Chunks"] { style: @ProductCard, icon: "ShieldCheck" }
      [agent_synth: "4. CitationSynthesizer" | "Grounded SSE Stream"] { style: @ProductCard, icon: "Gemini" }
    }
  }

  Zone @Path_B_Arena_And_Eval {
    title: "2B. Standard RAG, Model Arena & Autorater Judge"
    subtitle: "Flash vs. Pro Benchmark & LLM-as-a-Judge (/5)"
    style: @ZoneBlue, area: "archB"
    layout: row, gap: 22, padding: 18, align: "stretch", icon: "VertexAI", iconSize: 20

    Zone @MultiModel_Arena {
      title: "Multi-Model Arena", layout: column, gap: 16, padding: 14, align: "center",
      fill: "#ffffff", stroke: "#5b9bf3", cornerRadius: 10, icon: "Gemini", iconSize: 18
      [gemini_flash: "Gemini 3.5 / 3.8 Flash" | "Sub-400ms TTFT"] { style: @ProductCard, icon: "Gemini" }
      [gemini_pro: "Gemini 3.1 Pro" | "Deep Frontier Reasoning"] { style: @ProductCard, icon: "VertexAI" }
    }

    Zone @Autorater_And_Finops {
      title: "Quality & FinOps Audit", layout: column, gap: 16, padding: 14, align: "center",
      fill: "#ffffff", stroke: "#5b9bf3", cornerRadius: 10, icon: "Activity", iconSize: 18
      [vertex_autorater: "Vertex AI Autorater" | "Groundedness & Relevance /5"] { style: @ProductCard, icon: "SecurityCommandCenter" }
      [finops_badge: "Real-Time FinOps" | "Token Cost (~$0.00011)"] { style: @ProductCard, icon: "CloudBilling" }
    }
  }

  Zone @Shared_Data_Foundation {
    title: "3. Shared Hybrid Vector Index & Cloud Storage Persistence (Direct VPC Egress)"
    subtitle: "In-Memory 768-dim Cosine + BM25 Index Backed by Cloud Storage UBLA"
    style: @ZoneGreen, area: "store"
    layout: row, gap: 20, padding: 18, align: "center", icon: "GcpDatabase", iconSize: 20
    [gcs_corpus_bucket: "Cloud Storage UBLA" | "Corpus & JSONL Index"] { style: @ProductCard, icon: "GcpStorageBucket" }
    [vertex_embeddings: "Vertex Embeddings" | "text-embedding-004 (768d)"] { style: @ProductCard, icon: "AIPlatform" }
    [hybrid_engine: "Go Hybrid Engine" | "Cosine (0.7) + BM25 (0.3)"] { style: @ProductCard, icon: "Cpu" }
    [cloud_logging_eval: "Cloud Logging & Trace" | "SSE & Eval Telemetry"] { style: @ProductCard, icon: "CloudLogging" }
  }
}

[armor_iap] -> [cloud_run_go] { sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step", sequenceBadge: "1", badgeFill: $GcpBlue, badgeFontColor: $SurfaceWhite }
[cloud_run_go] -> [doc_ingestor] { sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[cloud_run_go] -> [agent_planner] { label: "Agentic SSE", sourceAnchor: "right", targetAnchor: "left", color: $GcpBlue, strokeWidth: 2, curve: "step", sequenceBadge: "2", badgeFill: $GcpBlue, badgeFontColor: $SurfaceWhite }
[agent_planner] -> [agent_retriever] { label: "Sub-Queries", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[agent_retriever] -> [agent_grader] { label: "Top-K Chunks", sourceAnchor: "right", targetAnchor: "left", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[agent_grader] -> [agent_synth] { label: "Verified", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step" }

[agent_synth] -> [gemini_flash] { label: "Compare", sourceAnchor: "right", targetAnchor: "left", color: $GcpBlue, strokeWidth: 1.8, dashed: true, curve: "step" }
[gemini_flash] -> [gemini_pro] { label: "Arena", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 1.8, curve: "step" }
[gemini_flash] -> [vertex_autorater] { label: "Judge /5", sourceAnchor: "right", targetAnchor: "left", color: $GcpGreen, strokeWidth: 2, curve: "step", sequenceBadge: "3", badgeFill: $GcpGreen, badgeFontColor: $SurfaceWhite }
[vertex_autorater] -> [finops_badge] { label: "Token Cost", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpGreen, strokeWidth: 1.8, curve: "step" }

[doc_ingestor] -> [gcs_corpus_bucket] { label: "Persist JSONL", sourceAnchor: "bottom", targetAnchor: "left", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[agent_retriever] -> [vertex_embeddings] { label: "768d Vector", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[agent_grader] -> [hybrid_engine] { label: "BM25 + Cos", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpBlue, strokeWidth: 2, curve: "step" }
[finops_badge] -> [cloud_logging_eval] { label: "Audit Logs", sourceAnchor: "bottom", targetAnchor: "top", color: $GcpGreen, strokeWidth: 1.8, dashed: true, curve: "step" }
```

---

## 2. Spécification GCP Draw (`go/gcpdraw`)

### 📋 Instructions d'utilisation
1. Rendez-vous sur l'outil officiel Google Cloud : **[GCP Draw (go/gcpdraw)](https://gcpdraw.corp.google.com)**.
2. Cliquez sur **Import / Code**.
3. Copiez-collez l'intégralité du bloc ci-dessous pour visualiser, éditer et exporter le schéma.

```text
meta {
  title "RAG Comparison Platform — 4-Engine Benchmark & CRAG Self-Correction (v3.3)"
}

elements {
  card users as users {
    display_name "AI Engineers & Decision Makers"
  }

  gcp {
    group frontend_ux {
      name "Interactive Benchmark UI (5 Modes + GCP X-Ray)"

      card run as web_ui {
        name "5 Modes & 1-Click Demo Pills"
        description "Standard, Split, Arena, Triple 3-Voies, Evaluation Dashboard"
      }
    }

    group cloud_run_backend {
      name "Cloud Run Gen2 Serverless Engine (Go 1.24 + RWMutex)"

      card run as rag_router {
        name "Concurrent SSE Fan-Out Router"
        description "Sub-millisecond TTFT Timer & FinOps Cost Estimator"
      }

      card run as engines_4 {
        name "4 Parallel Search Engines"
        description "1. Classical TF | 2. Dense Cosine 768d | 3. Hybrid RRF (k=60) | 4. Agentic CRAG"
      }

      card run as llm_judge {
        name "LLM-as-a-Judge & FinOps Telemetry"
        description "Faithfulness, Relevance, Completeness (/5) + Live Cost ($)"
      }
    }

    group vertex_ai_data {
      name "Vertex AI & Persistent Corpus Storage (europe-west1)"

      card vertex_ai as embeddings {
        name "Vertex AI text-embedding-004"
        description "768d Semantic Vectors & Batch Rate-Limit Shield"
      }

      card vertex_ai as gemini_flash {
        name "Vertex AI Gemini 2.5 Flash"
        description "Streaming Generation, CRAG Query Rewrite & Judge"
      }

      card storage as gcs_Persistence {
        name "Cloud Storage (GCS)"
        description "index.json Vector Persistence & Document Vault"
      }
    }
  }
}

paths {
  users -down-> web_ui : "1-Click Demo Pill / Upload PDF"
  web_ui --> rag_router : "SSE Stream (/api/query/stream)"
  rag_router --> engines_4 : "Parallel Execution"
  engines_4 --> embeddings : "768d Query Embedding"
  engines_4 --> gemini_flash : "CRAG Self-Correction & Generation"
  rag_router --> llm_judge : "Post-Stream Evaluation (/api/evaluate)"
  llm_judge --> gemini_flash : "Structured JSON Scoring"
  rag_router --> gcs_Persistence : "Auto-Save & Startup Hydrate"
}
```
