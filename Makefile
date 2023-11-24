.PHONY: slides

slides:
	hugo server --source slides/

slides_pdf:
	npx decktape@3 http://localhost:1313 slides.pdf
