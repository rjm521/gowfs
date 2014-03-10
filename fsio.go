package gowfs

import "fmt"
import "os"
import "io"
import "strconv"
import "strings"
import "net/url"
import "net/http"

// Creates a new file and stores its content in HDFS. 
// See HDFS FileSystem.create() 
// For detail, http://hadoop.apache.org/docs/stable/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Create_and_Write_to_a_File
// See NOTE section on that page for impl detail.
func (fs *FileSystem) Create(
	data io.Reader,
	p Path,
	overwrite bool, 
	blocksize uint64, 
	replication uint16, 
	permission os.FileMode, 
	buffersize uint) (bool, error){

	params := map[string]string{"op":OP_CREATE}
	params["overwrite"] = strconv.FormatBool(overwrite)
	
	if blocksize == 0 {
		params["blocksize"] = "134217728" // from hdfs-default.xml (ver 2)
	}else{
		params["blocksize"] = strconv.FormatInt(int64(blocksize), 10)
	}

	if replication == 0 {
		params["replication"] = "3"
	}else{
		params["replication"] = strconv.FormatInt(int64(replication), 10)
	}

	if permission <= 0 || permission > 1777 {
		params["permission"] = "0700"
	}else{
		params["permission"] = strconv.FormatInt(int64(permission), 8)
	}

	if buffersize == 0 {
		params["buffersize"] = "4096"
	}else{
		params["buffersize"] = strconv.FormatInt(int64(buffersize), 10)
	}

	u, err := buildRequestUrl(fs.Config, &p, &params)
	if err != nil {
		return false, err
	}

	// take over default transport to avoid redirect
	tr := &http.Transport{}
	req, _ := http.NewRequest("PUT", u.String(), nil)
	rsp, err := tr.RoundTrip(req)
	if err != nil {
		return false, err
	}

	// extract returned url in header.
	loc := rsp.Header.Get("Location")
	u, err = url.ParseRequestURI(loc)
	if err != nil {
		return false, fmt.Errorf("Create() - did not get a valid redirect URL from server.")
	}

	req,   _ = http.NewRequest("PUT", u.String(), data) 
	rsp, err = fs.client.Do(req)
	if  err != nil  {
		fmt.Errorf("Create() - bad url %s", loc)
		return false, err
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("Create() - File not created.  Server returned status %v", rsp.StatusCode)
	}

	return true, nil
}

//Opens the specificed Path and returns its content to be accessed locally. 
//See HDFS WebHdfsFileSystem.open()
// See http://hadoop.apache.org/docs/r2.2.0/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#HTTP_Query_Parameter_Dictionary
func (fs *FileSystem) Open(p Path, offset, length int64, buffSize int) (io.ReadCloser, error){
	params := map[string]string{"op":OP_OPEN}

	if offset < 0{
		params["offset"] = "0"
	}else{
		params["offset"] = strconv.FormatInt(offset, 10)
	}

	if length > 0{
		params["length"] = strconv.FormatInt(length, 10)
	}

	if buffSize <= 0 {
		params["buffersize"] = "4096"
	}else{
		params["buffersize"] = strconv.Itoa(buffSize)
	}

	u, err := buildRequestUrl(fs.Config, &p, &params)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", u.String(), nil)
	rsp, err := fs.client.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp.Body, nil
}

// Appends specified data to an existing file.
// See HDFS FileSystem.append()
// See http://hadoop.apache.org/docs/stable/hadoop-project-dist/hadoop-hdfs/WebHDFS.html#Append_to_a_File
func (fs *FileSystem) Append(data io.Reader, p Path, buffersize int)(bool, error){
	params := map[string]string{"op":OP_APPEND}
	
	if buffersize == 0 {
		params["buffersize"] = "4096"
	}else{
		params["buffersize"] = strconv.FormatInt(int64(buffersize), 10)
	}

	u, err := buildRequestUrl(fs.Config, &p, &params)
	if err != nil {
		return false, err
	}

	// take over default transport to avoid redirect
	tr := &http.Transport{}
	req, _ := http.NewRequest("POST", u.String(), nil)
	rsp, err := tr.RoundTrip(req)
	if err != nil {
		return false, err
	}

	// extract returned url in header.
	loc := rsp.Header.Get("Location")
	u, err = url.ParseRequestURI(loc)
	if err != nil {
		return false, fmt.Errorf("Append() - did not receive a valid URL from server.")
	}

	req,   _ = http.NewRequest("POST", u.String(), data) 
	rsp, err = fs.client.Do(req)
	if  err != nil  {
		return false, err
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("Create() - File not created.  Server returned status %v", rsp.StatusCode)
	}
	return true, nil
}

// Concatenate (on the server) a list of given files paths to a new file.
// See HDFS FileSystem.concat()
func (fs *FileSystem) Concat(target Path, sources []string)(bool, error) {
	if (target == Path{}) {
		return false, fmt.Errorf("Concat() - The target path must be provided.")
	}
	params := map[string]string{"op":OP_CONCAT}
	params["sources"] = strings.Join (sources, ",")
	
	u, err := buildRequestUrl(fs.Config, &target, &params)
	if err != nil {
		return false, err
	}

	req, _ 	 := http.NewRequest("POST", u.String(), nil)
	rsp, err := fs.client.Do(req)
	if err != nil {
		return false, err
	}
	if rsp.StatusCode != http.StatusOK && rsp.ContentLength != 0 {
		return false, fmt.Errorf("Concat() - Server returned unexpected result.")
	}
	return true, nil
}