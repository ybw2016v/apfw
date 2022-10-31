addEventListener("fetch", (event) => {
    event.respondWith(
      handleRequest(event.request).catch(
        (err) => new Response(err.stack, { status: 500 })
      )
    );
  });
  
  /**
   * Many more examples available at:
   *   https://developers.cloudflare.com/workers/examples
   * @param {Request} request
   * @returns {Promise<Response>}
   */
  
  async function handleRequest(request) {
    const url = new URL(request.url);
    const oldheader = request.headers;
    const tkey = oldheader.get("tkey");
    const ua = oldheader.get("user-agent")
    if (tkey!=TOKEN){
      return new Response("503 Bad Gateway",{status:403})
    }
    const thost=oldheader.get("thost");
    const nhd = new Headers();
    Object.assign(nhd,oldheader);
    nhd.set('host',thost);
    nhd.set('user-agent',ua+" (Forwarded by apfw-cf v0.1)")
    nhd.delete('tkey');
    nhd.delete('thost');
    url.hostname=thost; 
    res = await fetch(url,{method:request.method,headers:nhd,body:request.body});
    const restp=res.headers.get("Content-Type")
    return new Response(res.body, {
        headers: { "Content-Type":restp },status:res.status
      });
  }