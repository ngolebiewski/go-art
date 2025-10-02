import { ARTWORK_DATA } from "../data/mockArtworks"

const Gallery = () => {
  const artworks = ARTWORK_DATA;

  return (
    <section>
      <h2>Gallery (refresh on map)</h2>
      <div id="art-gallery">

        {artworks.map(artwork => ( 
          <article key={artwork.id}>
            <h3>{artwork.title}</h3>
            <img src={artwork.thumbnailUrl} alt={artwork.title}/>
            <p>{artwork.ownerName}</p>
          </article>
        )
        )}

      </div>
    </section>
  )
}
export default Gallery