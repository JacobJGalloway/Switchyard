using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Switchyard.InventoryAPI.Data.Migrations
{
    /// <inheritdoc />
    public partial class AddUnitPrice : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<string>(
                name: "PriceCurrency",
                table: "Tools",
                type: "text",
                nullable: false,
                defaultValue: "USD");

            migrationBuilder.AddColumn<DateTime>(
                name: "PriceEffectiveDate",
                table: "Tools",
                type: "timestamp with time zone",
                nullable: true);

            migrationBuilder.AddColumn<decimal>(
                name: "UnitPrice",
                table: "Tools",
                type: "numeric",
                nullable: false,
                defaultValue: 0m);

            migrationBuilder.AddColumn<string>(
                name: "PriceCurrency",
                table: "PPE",
                type: "text",
                nullable: false,
                defaultValue: "USD");

            migrationBuilder.AddColumn<DateTime>(
                name: "PriceEffectiveDate",
                table: "PPE",
                type: "timestamp with time zone",
                nullable: true);

            migrationBuilder.AddColumn<decimal>(
                name: "UnitPrice",
                table: "PPE",
                type: "numeric",
                nullable: false,
                defaultValue: 0m);

            migrationBuilder.AddColumn<string>(
                name: "PriceCurrency",
                table: "Clothing",
                type: "text",
                nullable: false,
                defaultValue: "USD");

            migrationBuilder.AddColumn<DateTime>(
                name: "PriceEffectiveDate",
                table: "Clothing",
                type: "timestamp with time zone",
                nullable: true);

            migrationBuilder.AddColumn<decimal>(
                name: "UnitPrice",
                table: "Clothing",
                type: "numeric",
                nullable: false,
                defaultValue: 0m);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "PriceCurrency",
                table: "Tools");

            migrationBuilder.DropColumn(
                name: "PriceEffectiveDate",
                table: "Tools");

            migrationBuilder.DropColumn(
                name: "UnitPrice",
                table: "Tools");

            migrationBuilder.DropColumn(
                name: "PriceCurrency",
                table: "PPE");

            migrationBuilder.DropColumn(
                name: "PriceEffectiveDate",
                table: "PPE");

            migrationBuilder.DropColumn(
                name: "UnitPrice",
                table: "PPE");

            migrationBuilder.DropColumn(
                name: "PriceCurrency",
                table: "Clothing");

            migrationBuilder.DropColumn(
                name: "PriceEffectiveDate",
                table: "Clothing");

            migrationBuilder.DropColumn(
                name: "UnitPrice",
                table: "Clothing");
        }
    }
}
