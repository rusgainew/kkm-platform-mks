import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Card, CardHeader, CardContent, CardFooter } from "../Card";

describe("Card", () => {
  it("should render card with content", () => {
    render(<Card>Card Content</Card>);
    expect(screen.getByText("Card Content")).toBeInTheDocument();
  });

  it("should apply default variant", () => {
    const { container } = render(<Card>Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("bg-gray-900");
    expect(card.className).toContain("border-gray-800");
  });

  it("should apply bordered variant", () => {
    const { container } = render(<Card variant="bordered">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("border-gray-700");
  });

  it("should apply elevated variant", () => {
    const { container } = render(<Card variant="elevated">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("shadow-lg");
  });

  it("should apply small padding", () => {
    const { container } = render(<Card padding="sm">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("p-4");
  });

  it("should apply large padding", () => {
    const { container } = render(<Card padding="lg">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("p-8");
  });

  it("should apply no padding", () => {
    const { container } = render(<Card padding="none">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).not.toContain("p-");
  });

  it("should apply hover effect when hover prop is true", () => {
    const { container } = render(<Card hover>Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("hover:shadow-lg");
  });

  it("should apply custom className", () => {
    const { container } = render(<Card className="custom-class">Content</Card>);
    const card = container.firstChild as HTMLElement;
    expect(card.className).toContain("custom-class");
  });
});

describe("CardHeader", () => {
  it("should render title as string", () => {
    render(<CardHeader title="Card Title" />);
    expect(screen.getByText("Card Title")).toBeInTheDocument();
  });

  it("should render title as ReactNode", () => {
    render(
      <CardHeader
        title={
          <span>
            Custom <strong>Title</strong>
          </span>
        }
      />
    );
    expect(screen.getByText("Custom")).toBeInTheDocument();
    expect(screen.getByText("Title")).toBeInTheDocument();
  });

  it("should render subtitle", () => {
    render(<CardHeader title="Title" subtitle="Subtitle" />);
    expect(screen.getByText("Subtitle")).toBeInTheDocument();
  });

  it("should render action element", () => {
    render(<CardHeader title="Title" action={<button>Action</button>} />);
    expect(screen.getByText("Action")).toBeInTheDocument();
  });

  it("should render children", () => {
    render(<CardHeader>Custom Content</CardHeader>);
    expect(screen.getByText("Custom Content")).toBeInTheDocument();
  });
});

describe("CardContent", () => {
  it("should render content", () => {
    render(<CardContent>Content Area</CardContent>);
    expect(screen.getByText("Content Area")).toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(
      <CardContent className="custom-content">Content</CardContent>
    );
    expect(container.firstChild).toHaveClass("custom-content");
  });
});

describe("CardFooter", () => {
  it("should render footer content", () => {
    render(<CardFooter>Footer Content</CardFooter>);
    expect(screen.getByText("Footer Content")).toBeInTheDocument();
  });

  it("should apply right alignment by default", () => {
    const { container } = render(<CardFooter>Footer</CardFooter>);
    const footer = container.firstChild as HTMLElement;
    expect(footer.className).toContain("justify-end");
  });

  it("should apply left alignment", () => {
    const { container } = render(<CardFooter align="left">Footer</CardFooter>);
    const footer = container.firstChild as HTMLElement;
    expect(footer.className).toContain("justify-start");
  });

  it("should apply center alignment", () => {
    const { container } = render(<CardFooter align="center">Footer</CardFooter>);
    const footer = container.firstChild as HTMLElement;
    expect(footer.className).toContain("justify-center");
  });

  it("should apply between alignment", () => {
    const { container } = render(<CardFooter align="between">Footer</CardFooter>);
    const footer = container.firstChild as HTMLElement;
    expect(footer.className).toContain("justify-between");
  });
});

describe("Card Composition", () => {
  it("should render complete card with all parts", () => {
    render(
      <Card>
        <CardHeader title="Test Card" subtitle="A test subtitle" />
        <CardContent>Main content area</CardContent>
        <CardFooter>
          <button>Cancel</button>
          <button>Submit</button>
        </CardFooter>
      </Card>
    );

    expect(screen.getByText("Test Card")).toBeInTheDocument();
    expect(screen.getByText("A test subtitle")).toBeInTheDocument();
    expect(screen.getByText("Main content area")).toBeInTheDocument();
    expect(screen.getByText("Cancel")).toBeInTheDocument();
    expect(screen.getByText("Submit")).toBeInTheDocument();
  });
});
