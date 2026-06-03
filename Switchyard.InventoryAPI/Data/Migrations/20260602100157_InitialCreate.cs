using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Switchyard.InventoryAPI.Data.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Clothing",
                columns: table => new
                {
                    PartitionKey = table.Column<string>(type: "text", nullable: false),
                    RowKey = table.Column<string>(type: "text", nullable: false),
                    LocationId = table.Column<string>(type: "text", nullable: false),
                    SKUMarker = table.Column<string>(type: "text", nullable: false),
                    UnloadedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    Projected = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Clothing", x => x.PartitionKey);
                });

            migrationBuilder.CreateTable(
                name: "PPE",
                columns: table => new
                {
                    PartitionKey = table.Column<string>(type: "text", nullable: false),
                    RowKey = table.Column<string>(type: "text", nullable: false),
                    LocationId = table.Column<string>(type: "text", nullable: false),
                    SKUMarker = table.Column<string>(type: "text", nullable: false),
                    UnloadedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    Projected = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_PPE", x => x.PartitionKey);
                });

            migrationBuilder.CreateTable(
                name: "Tools",
                columns: table => new
                {
                    PartitionKey = table.Column<string>(type: "text", nullable: false),
                    RowKey = table.Column<string>(type: "text", nullable: false),
                    LocationId = table.Column<string>(type: "text", nullable: false),
                    SKUMarker = table.Column<string>(type: "text", nullable: false),
                    UnloadedDate = table.Column<DateTime>(type: "timestamp with time zone", nullable: false),
                    Projected = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Tools", x => x.PartitionKey);
                });
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "Clothing");

            migrationBuilder.DropTable(
                name: "PPE");

            migrationBuilder.DropTable(
                name: "Tools");
        }
    }
}
