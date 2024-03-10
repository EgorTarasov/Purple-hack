import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export default function ModelSelect() {
  return (
    <Select>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Модель" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Модель</SelectLabel>
          <SelectItem value="apple">LLAMA 2</SelectItem>
          <SelectItem value="banana">FRED T5</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
