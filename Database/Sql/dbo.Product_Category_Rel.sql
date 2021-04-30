CREATE TABLE [dbo].[Product_Category_Rel]
(
	[Id] BIGINT NOT NULL  IDENTITY, 
    [ProductId] BIGINT NOT NULL, 
    [CategoryId] BIGINT NOT NULL, 
    PRIMARY KEY ([ProductId], [CategoryId])
)
