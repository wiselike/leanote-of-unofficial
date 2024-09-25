all: build

.PHONY:fmt build release

# format golang code
fmt:
	find -name "*.go" -exec go fmt {} \;

# only build temporarily
build:
	rm -rf tmp/
	GOPATH=`go env GOPATH` revel build . tmp/

# build js
gulp:
	@cd public/tinymce; rm -f tinymce.js tinymce.dev.js tinymce.min.js tinymce.jquery.dev.js tinymce.full.js tinymce.full.min.js
	@cd public/tinymce; grunt minify;
	@cd public/tinymce; grunt bundle --themes=leanote --plugins=autolink,link,leaui_image,leaui_mindmap,lists,hr,paste,searchreplace,leanote_nav,leanote_code,tabfocus,table,directionality,textcolor;
	@gulp

# build all and rerun leanote
release: gulp
	GOPATH=`go env GOPATH` revel build . release/
	rsync -azr --delete --delete-before --exclude github.com/wiselike/leanote-of-unofficial/conf/app.conf --exclude github.com/wiselike/leanote-of-unofficial/public/upload --exclude github.com/wiselike/leanote-of-unofficial/mongodb_backup -e 'ssh -p 22' release/src/ root@192.168.0.12:/root/dockers/leanote/leanote/src
	rsync -azr release/leanote-of-unofficial  -e 'ssh -p 22' root@192.168.0.12:/root/dockers/leanote/leanote/leanote-of-unofficial
	rm -rf release/
	ssh -p 22 root@192.168.0.12 "docker restart leanote"
