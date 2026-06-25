using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Switchyard.LogisticsAPI.Data.Migrations
{
    /// <inheritdoc />
    public partial class AddCommittedDate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<DateTime>(
                name: "CommittedDate",
                table: "BillsOfLading",
                type: "timestamp with time zone",
                nullable: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "CommittedDate",
                table: "BillsOfLading");
        }
    }
}
