# http://revel.github.io/manual/tool.html

# go get -u github.com/wiselike/revel-cmd/revel
SCRIPTPATH=$(dirname "$PWD")
echo $SCRIPTPATH;
cd $SCRIPTPATH;

revel run -a .