<?php


define( 'CSAJAX_FILTERS', true );
define( 'CSAJAX_FILTER_DOMAIN', true );
define( 'CSAJAX_DEBUG', true );
$request_headers = array( );
$setContentType = true;
$isMultiPart = false;
foreach ( $_SERVER as $key => $value ) {
    if(preg_match('/Content.Type/i', $key)){
        $setContentType = false;
        $content_type = explode(";", $value)[0];
        $isMultiPart = preg_match('/multipart/i', $content_type);
        $request_headers[] = "Content-Type: ".$content_type;
        continue;
    }
  if ( substr( $key, 0, 5 ) == 'HTTP_' ) {
    $headername = str_replace( '_', ' ', substr( $key, 5 ) );
    $headername = str_replace( ' ', '-', ucwords( strtolower( $headername ) ) );
    if ( !in_array( $headername, array( 'Host', 'X-Proxy-Url' ) ) ) {
      $request_headers[] = "$headername: $value";
    }
  }
}

if($setContentType)
    $request_headers[] = "Content-Type: application/json";

// identify request method, url and params
$request_method = $_SERVER['REQUEST_METHOD'];
if ( 'GET' == $request_method ) {
  $request_params = $_GET;
} elseif ( 'POST' == $request_method ) {
  $request_params = $_POST;
  if ( empty( $request_params ) ) {
    $data = file_get_contents( 'php://input' );
    if ( !empty( $data ) ) {
      $request_params = $data;
    }
  }
} elseif ( 'PUT' == $request_method || 'DELETE' == $request_method ) {
  $request_params = file_get_contents( 'php://input' );
} else {
  $request_params = null;
}

// Get URL from `csurl` in GET or POST data, before falling back to X-Proxy-URL header.
if ( isset( $_REQUEST['csurl'] ) ) {
    $request_url = $_REQUEST['csurl'];
} else {
    header( $_SERVER['SERVER_PROTOCOL'] . ' 404 Not Found');
    header( 'Status: 404 Not Found' );
    $_SERVER['REDIRECT_STATUS'] = 404;
    exit;
}

$p_request_url = parse_url( $request_url );

if($p_request_url['scheme'] == 'http' && $p_request_url['host'] == '127.0.0.1') {
  if(!array_key_exists('port', $p_request_url)) {
    $p_request_url['port'] = '80';
  }
  if($p_request_url['port'] == '80' && ($p_request_url['path'] == '/' || $p_request_url['path'] == '')) {
    $p_request_url['path'] = '/InternalServices.php';
    $request_url = $p_request_url['scheme'] . "://" . $p_request_url['host'] . ":" . $p_request_url['port'] . $p_request_url['path'];
  }

}

if(strpos($request_url, 'cron') !== false) {
  echo "BLOCKED BY WAF!";
  return;
}
// csurl may exist in GET request methods
if ( is_array( $request_params ) && array_key_exists('csurl', $request_params ) )
  unset( $request_params['csurl'] );

// ignore requests for proxy :)
if ( preg_match( '!' . $_SERVER['SCRIPT_NAME'] . '!', $request_url ) || empty( $request_url ) || count( $p_request_url ) == 1 ) {
  // csajax_debug_message( 'Invalid request - make sure that csurl variable is not empty' );
  exit;
}

// check against valid requests
if ( CSAJAX_FILTERS ) {
  $parsed = $p_request_url;
  $check_url = isset( $parsed['scheme'] ) ? $parsed['scheme'] . '://' : '';
  $check_url .= isset( $parsed['user'] ) ? $parsed['user'] . ($parsed['pass'] ? ':' . $parsed['pass'] : '') . '@' : '';
  $check_url .= isset( $parsed['host'] ) ? $parsed['host'] : '';
  $check_url .= isset( $parsed['port'] ) ? ':' . $parsed['port'] : '';
  $check_url .= isset( $parsed['path'] ) ? $parsed['path'] : '';
  
}

// append query string for GET requests
if ( $request_method == 'GET' && count( $request_params ) > 0 && (!array_key_exists( 'query', $p_request_url ) || empty( $p_request_url['query'] ) ) ) {
  echo $request_params;
  exit();
  $request_url .= '?' . http_build_query( $request_params );
}


// let the request begin

$ch = curl_init( $request_url );
curl_setopt( $ch, CURLOPT_HTTPHEADER, $request_headers );   // (re-)send headers
curl_setopt( $ch, CURLOPT_RETURNTRANSFER, true );  // return response
curl_setopt( $ch, CURLOPT_HEADER, true );    // enabled response headers
// add data for POST, PUT or DELETE requests
if ( 'POST' == $request_method ) {
  $post_data = is_array( $request_params ) ? http_build_query( $request_params ) : $request_params;

    $has_files = false;
    $file_params = array();

    foreach ($_FILES as $f => $file) {
        if($file['size']){
            $file_params[$f] = '@'. $file['tmp_name'] .";type=". $file['type'];
            $has_files = true;
        }
    }

    if($isMultiPart || $has_files){
        foreach(explode("&",$post_data) as $i => $param) {
            $params = explode("=", $param);
            $xvarname = $params[0];
            if (!empty($xvarname))
                $file_params[$xvarname] = $params[1];
        }
    }

  curl_setopt( $ch, CURLOPT_POST, true );
  curl_setopt( $ch, CURLOPT_POSTFIELDS,  $isMultiPart || $has_files ? $file_params : $post_data );
} elseif ( 'PUT' == $request_method || 'DELETE' == $request_method ) {
  curl_setopt( $ch, CURLOPT_CUSTOMREQUEST, $request_method );
  curl_setopt( $ch, CURLOPT_POSTFIELDS, $request_params );
}

// retrieve response (headers and content)
$response = curl_exec( $ch );
curl_close( $ch );

// split response to header and content
list($response_headers, $response_content) = preg_split( '/(\r\n){2}/', $response, 2 );

// (re-)send the headers
$response_headers = preg_split( '/(\r\n){1}/', $response_headers );
foreach ( $response_headers as $key => $response_header ) {
  // Rewrite the `Location` header, so clients will also use the proxy for redirects.
  if ( preg_match( '/^Location:/', $response_header ) ) {
    list($header, $value) = preg_split( '/: /', $response_header, 2 );
    $response_header = 'Location: http://' . $_SERVER['HTTP_HOST'] . '/?csurl=' . $value;
  }
  if ( !preg_match( '/^(Transfer-Encoding):/', $response_header ) ) {
    header( $response_header, false );
  }
}


// // finally, output the content

print($response_content);
