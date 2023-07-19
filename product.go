package shopify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gempages/go-shopify-graphql/graphql"
	"github.com/gempages/go-shopify-graphql/utils"

	"github.com/r0busta/go-shopify-graphql-model/v3/graph/model"
)

type ProductService interface {
	List(query string) ([]*model.Product, error)
	ListAll() ([]*model.Product, error)
	ListWithFields(query string, fields string, first int, after string) (*model.ProductConnection, error)

	Get(id string) (*model.Product, error)
	GetWithFields(id string, fields string) (*model.Product, error)
	GetSingleProductCollection(id string, cursor string) (*model.Product, error)
	GetSingleProductVariant(id string, cursor string) (*model.Product, error)
	GetSingleProduct(id string) (*model.Product, error)

	Create(product model.ProductInput, media []model.CreateMediaInput) error
	Update(product model.ProductInput) error
	Delete(product model.ProductDeleteInput) error
}

type ProductServiceOp struct {
	client *Client
}

var _ ProductService = &ProductServiceOp{}

type ProductBase struct {
	ID               graphql.ID           `json:"id,omitempty"`
	CreatedAt        time.Time            `json:"createdAt,omitempty"`
	LegacyResourceID graphql.String       `json:"legacyResourceId,omitempty"`
	Handle           graphql.String       `json:"handle,omitempty"`
	Options          []ProductOption      `json:"options,omitempty"`
	Tags             []graphql.String     `json:"tags,omitempty"`
	Description      graphql.String       `json:"description,omitempty"`
	Title            graphql.String       `json:"title,omitempty"`
	PriceRangeV2     *ProductPriceRangeV2 `json:"priceRangeV2,omitempty"`
	ProductType      graphql.String       `json:"productType,omitempty"`
	Vendor           graphql.String       `json:"vendor,omitempty"`
	TotalInventory   graphql.Int          `json:"totalInventory,omitempty"`
	OnlineStoreURL   graphql.String       `json:"onlineStoreUrl,omitempty"`
	DescriptionHTML  graphql.String       `json:"descriptionHtml,omitempty"`
	SEO              Seo                  `json:"seo,omitempty"`
	TemplateSuffix   graphql.String       `json:"templateSuffix,omitempty"`
	Status           graphql.String       `json:"status,omitempty"`
	PublishedAt      *time.Time           `json:"publishedAt,omitempty"`
	UpdatedAt        time.Time            `json:"updatedAt,omitempty"`
	TracksInventory  bool                 `json:"tracksInventory,omitempty"`
}

type ProductBulkResult struct {
	ProductBase

	Metafields      []Metafield      `json:"metafields,omitempty"`
	ProductVariants []ProductVariant `json:"variants,omitempty"`
	Collections     []Collection     `json:"collections,omitempty"`
	ProductImages   []ProductImage   `json:"images,omitempty"`
	Media           []Media          `json:"media,omitempty"`
}

type ProductImage struct {
	AltText graphql.String `json:"altText,omitempty"`
	ID      graphql.ID     `json:"id,omitempty"`
	Src     graphql.String `json:"src,omitempty"`
	Height  graphql.Int    `json:"height,omitempty"`
	Width   graphql.Int    `json:"width,omitempty"`
}

type Media struct {
	ID               graphql.ID             `json:"id,omitempty"`
	MimeType         graphql.String         `json:"mimeType,omitempty"`
	MediaContentType model.MediaContentType `json:"mediaContentType,omitempty"`
	Alt              graphql.String         `json:"alt,omitempty"`
	Image            *ProductImage          `json:"image,omitempty"`
	Sources          interface{}            `json:"sources,omitEmpty"`
	OriginalSource   *Source                `json:"originalSource,omitempty"`
	EmbedUrl         graphql.String         `json:"embedUrl,omitempty"`
	OriginUrl        graphql.String         `json:"originUrl,omitempty"`
	Preview          Preview                `json:"preview,omitempty"`
}

type Preview struct {
	Image ProductImage `json:"image,omitempty"`
}

type Source struct {
	MimeType graphql.String `json:"mimeType,omitempty"`
	Url      graphql.String `json:"url,omitempty"`
	FileSize graphql.Int    `json:"filesize,omitempty"`
	Format   graphql.String `json:"format,omitempty"`
}

// SEO information.
type Seo struct {
	// SEO Description.
	Description graphql.String `json:"description,omitempty"`
	// SEO Title.
	Title graphql.String `json:"title,omitempty"`
}

type ProductOption struct {
	ID       graphql.ID       `json:"id,omitempty"`
	Name     graphql.String   `json:"name,omitempty"`
	Position graphql.Int      `json:"position,omitempty"`
	Values   []graphql.String `json:"values,omitempty"`
}

type ProductPriceRangeV2 struct {
	MinVariantPrice MoneyV2 `json:"minVariantPrice,omitempty"`
	MaxVariantPrice MoneyV2 `json:"maxVariantPrice,omitempty"`
}

type ProductUpdate struct {
	ProductInput model.ProductInput
}

type MetafieldInput struct {
	ID        graphql.ID               `json:"id,omitempty"`
	Namespace graphql.String           `json:"namespace,omitempty"`
	Key       graphql.String           `json:"key,omitempty"`
	Value     graphql.String           `json:"value,omitempty"`
	Type      model.MetafieldValueType `json:"type,omitempty"`
}

type SEOInput struct {
	Description graphql.String `json:"description,omitempty"`
	Title       graphql.String `json:"title,omitempty"`
}

type ImageInput struct {
	AltText graphql.String `json:"altText,omitempty"`
	ID      graphql.ID     `json:"id,omitempty"`
	Src     graphql.String `json:"src,omitempty"`
}

type mutationProductCreate struct {
	ProductCreateResult productCreateResult `graphql:"productCreate(input: $input, media: $media)" json:"productCreate"`
}

type mutationProductUpdate struct {
	ProductUpdateResult productUpdateResult `graphql:"productUpdate(input: $input)" json:"productUpdate"`
}

type mutationProductDelete struct {
	ProductDeleteResult productDeleteResult `graphql:"productDelete(input: $input)" json:"productDelete"`
}

type productCreateResult struct {
	Product    *model.Product    `json:"product,omitempty"`
	UserErrors []model.UserError `json:"userErrors,omitempty"`
}

type productUpdateResult struct {
	Product    *model.Product `json:"product,omitempty"`
	UserErrors []UserErrors   `json:"userErrors"`
}

type productDeleteResult struct {
	ID         string       `json:"deletedProductId,omitempty"`
	UserErrors []UserErrors `json:"userErrors"`
}

const productBaseQuery = `
  id
  legacyResourceId
  handle
  status
  publishedAt
  createdAt
  updatedAt
  tracksInventory
	options{
    	id
		name
		position
		values
	}
	tags
	title
	description
	priceRangeV2{
		minVariantPrice{
			amount
			currencyCode
		}
		maxVariantPrice{
			amount
			currencyCode
		}
	}
	productType
	vendor
	totalInventory
	onlineStoreUrl
	descriptionHtml
	seo{
		description
		title
	}
	templateSuffix
`

var singleProductQueryVariant = fmt.Sprintf(`
  id
  variants(first: 100) {
    edges {
      node {
        id
		createdAt
		updatedAt
        legacyResourceId
        sku
        selectedOptions {
          name
          value
        }
        compareAtPrice
        price
        inventoryQuantity
		image {
			altText
			height
			id
			src
			width
		}
        barcode
        title
        inventoryPolicy
        inventoryManagement
        weightUnit
        weight
		position
      }
      cursor
    }
  }

`)

var singleProductQueryVariantWithCursor = fmt.Sprintf(`
  id
  variants(first: 100, after: $cursor) {
    edges {
      node {
        id
		createdAt
		updatedAt
        legacyResourceId
        sku
        selectedOptions {
          name
          value
        }
        compareAtPrice
        price
        inventoryQuantity
        barcode
        title
        inventoryPolicy
        inventoryManagement
        weightUnit
        weight
		position
      }
      cursor
    }
  }

`)

var singleProductQueryCollection = fmt.Sprintf(`
  id
  collections(first:250) {
    edges {
      node {
        id
        title
        handle
        description
        templateSuffix
       	image {
			altText
			height
			id
			src
			width
		}
      }
      cursor
    }
  }
`)

var singleProductQueryCollectionWithCursor = fmt.Sprintf(`
  id
  collections(first:250, after: $cursor) {
    edges {
      node {
        id
		title
        handle
        description
        templateSuffix
		image {
			altText
			height
			id
			src
			width
		}
      }
      cursor
    }
  }
`)

var productQuery = fmt.Sprintf(`
	%s
	variants(first:100, after: $cursor){
		edges{
			node{
				id
				createdAt
				updatedAt
				legacyResourceId
				sku
				selectedOptions{
					name
					value
				}
				compareAtPrice
				price
				inventoryQuantity
				barcode
				title
				inventoryPolicy
				inventoryManagement
				weightUnit
				weight
				position
			}
		}
		pageInfo{
			hasNextPage
		}
	}
`, productBaseQuery)

var productBulkQuery = fmt.Sprintf(`
	%s
	metafields{
		edges{
			node{
				id
				legacyResourceId
				namespace
				key
				value
				type
			}
		}
	}
    images {
        edges {
            node {
                altText
                height
                id
                src
                width
            }
        }
    }
	media {
		edges {
			node {
				mediaContentType
				...on MediaImage {
					id
					alt
					mimeType
					image {
                		height
                		src
                		width
					}
				}
				...on Model3d {
					id
					alt
					originalSource {
						url
						format
						filesize
						mimeType
					}
					preview {
						image {
							src
						}
					}
				}
				...on Video {
					id
					alt
					duration
					originalSource {
						url
						format
						mimeType
 						height
						width
					}
					preview {
						image {
							src
						}
					}
				}
				...on ExternalVideo {
					id
					originUrl
					embedUrl
					preview {
						image {
							src
						}
					}
				}
			}
		}
	}
	variants{
		edges{
			node{
				id
				createdAt
				updatedAt
				legacyResourceId
				sku
				selectedOptions{
					name
					value
				}
                image {
                    altText
                    height
                    id
                    src
                    width
                }
				compareAtPrice
				price
				inventoryQuantity
				barcode
				title
				inventoryPolicy
				inventoryManagement
				weightUnit
				weight
				position
			}
		}
	}
`, productBaseQuery)

func (s *ProductServiceOp) ListAll() ([]*model.Product, error) {
	q := fmt.Sprintf(`
		query products{
			products{
				edges{
					node{
						%s
					}
				}
			}
		}
	`, productBulkQuery)

	res := make([]*model.Product, 0)
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProductServiceOp) List(query string) ([]*model.Product, error) {
	q := fmt.Sprintf(`
		query products{
			products(query: "$query"){
				edges{
					node{
						%s
					}
				}
			}
		}
	`, productBulkQuery)

	q = strings.ReplaceAll(q, "$query", query)

	res := make([]*model.Product, 0)
	err := s.client.BulkOperation.BulkQuery(q, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProductServiceOp) Get(id string) (*model.Product, error) {
	out, err := s.getPage(id, "")
	if err != nil {
		return nil, err
	}

	nextPageData := out
	if out != nil && out.Variants != nil && out.Variants.PageInfo != nil {
		hasNextPage := out.Variants.PageInfo.HasNextPage
		for hasNextPage && len(nextPageData.Variants.Edges) > 0 {
			cursor := nextPageData.Variants.Edges[len(nextPageData.Variants.Edges)-1].Cursor
			nextPageData, err := s.getPage(id, cursor)
			if err != nil {
				return nil, err
			}
			out.Variants.Edges = append(out.Variants.Edges, nextPageData.Variants.Edges...)
			hasNextPage = nextPageData.Variants.PageInfo.HasNextPage
		}
	}

	return out, nil
}

func (s *ProductServiceOp) getPage(id string, cursor string) (*model.Product, error) {
	q := fmt.Sprintf(`
		query product($id: ID!, $cursor: String) {
			product(id: $id){
				%s
			}
		}
	`, productQuery)

	vars := map[string]interface{}{
		"id": id,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	out := model.QueryRoot{}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) GetWithFields(id string, fields string) (*model.Product, error) {
	if fields == "" {
		fields = `id`
	}
	q := fmt.Sprintf(`
		query product($id: ID!) {
		  product(id: $id){
			%s
		  }
		}`, fields)

	vars := map[string]interface{}{
		"id": id,
	}

	out := model.QueryRoot{}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) GetSingleProductCollection(id string, cursor string) (*model.Product, error) {
	q := ""
	if cursor != "" {
		q = fmt.Sprintf(`
    query product($id: ID!, $cursor: String) {
      product(id: $id){
        %s
      }
    }
    `, singleProductQueryCollectionWithCursor)
	} else {
		q = fmt.Sprintf(`
    query product($id: ID!) {
      product(id: $id){
        %s
      }
    }
    `, singleProductQueryCollection)
	}

	vars := map[string]interface{}{
		"id": id,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	out := model.QueryRoot{}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) GetSingleProductVariant(id string, cursor string) (*model.Product, error) {
	q := ""
	if cursor != "" {
		q = fmt.Sprintf(`
    query product($id: ID!, $cursor: String) {
      product(id: $id){
        %s
      }
    }
    `, singleProductQueryVariantWithCursor)
	} else {
		q = fmt.Sprintf(`
    query product($id: ID!) {
      product(id: $id){
        %s
      }
    }
    `, singleProductQueryVariant)
	}

	vars := map[string]interface{}{
		"id": id,
	}
	if cursor != "" {
		vars["cursor"] = cursor
	}

	out := model.QueryRoot{}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) GetSingleProduct(id string) (*model.Product, error) {
	q := fmt.Sprintf(`
		query product($id: ID!) {
			product(id: $id){
				%s
				%s
				%s
			}
		}
	`, productBaseQuery, singleProductQueryVariant, singleProductQueryCollection)

	vars := map[string]interface{}{
		"id": id,
	}

	out := model.QueryRoot{}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Product, nil
}

func (s *ProductServiceOp) ListWithFields(query, fields string, first int, after string) (*model.ProductConnection, error) {
	if fields == "" {
		fields = `id`
	}

	q := fmt.Sprintf(`
		query products ($first: Int!, $after: String, $query: String) {
			products (first: $first, after: $after, query: $query) {
				edges {
					cursor
					node {
						%s
					}
				}
				pageInfo {
					hasNextPage
				}
			}
		}
	`, fields)

	vars := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		vars["after"] = after
	}
	if query != "" {
		vars["query"] = query
	}
	out := model.QueryRoot{}

	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.QueryString(context.Background(), q, vars, &out)
	})
	if err != nil {
		return nil, err
	}

	return out.Products, nil
}

func (s *ProductServiceOp) Create(product model.ProductInput, media []model.CreateMediaInput) error {
	m := mutationProductCreate{}

	vars := map[string]interface{}{
		"input": product,
		"media": media,
	}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.Mutate(context.Background(), &m, vars)
	})
	if err != nil {
		return err
	}

	if len(m.ProductCreateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductCreateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) Update(product model.ProductInput) error {
	m := mutationProductUpdate{}

	vars := map[string]interface{}{
		"input": product,
	}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.Mutate(context.Background(), &m, vars)
	})
	if err != nil {
		return err
	}

	if len(m.ProductUpdateResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductUpdateResult.UserErrors)
	}

	return nil
}

func (s *ProductServiceOp) Delete(product model.ProductDeleteInput) error {
	m := mutationProductDelete{}

	vars := map[string]interface{}{
		"input": product,
	}
	err := utils.ExecWithRetries(s.client.retries, func() error {
		return s.client.gql.Mutate(context.Background(), &m, vars)
	})
	if err != nil {
		return err
	}

	if len(m.ProductDeleteResult.UserErrors) > 0 {
		return fmt.Errorf("%+v", m.ProductDeleteResult.UserErrors)
	}

	return nil
}
